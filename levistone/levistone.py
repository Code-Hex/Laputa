import nfc
import nfc.clf
import nfc.ndef
import binascii
from pprint import pprint
import sys
import socket
import pifacedigitalio
import threading
import time

import logging
logging.basicConfig(level=logging.INFO)
log = logging.getLogger('main')
handler = logging.FileHandler('scan_card.log')
handler.setLevel(logging.INFO)

formatter = logging.Formatter('%(levelname)s(%(name)s): %(asctime)s - %(message)s')
handler.setFormatter(formatter)
log.addHandler(handler)


class Reader:
	def __init__(self, socket_path):
		# unix domain socket
		self.socket_path = socket_path
		self.socket = socket.socket(socket.AF_UNIX,socket.SOCK_STREAM)

		# nfc library
		self.system_code = 0xFE00
		self.service_code = {
			"suica":  0x090f,
			"univ":   0x50CB,
			"edy":    0x110B,
			"waon":   0x67CF,
			"nanaco": 0x558B
		}
		
		# piface
		self.piface = pifacedigitalio.PiFaceDigital()

	def __del__(self):
		self.socket.close()

	def socket_connect(self):
		self.socket.connect(self.socket_path)
        
	def send(self, msg):
		log.info("send to server: %s" % msg)
		#self.socket.send(msg.encode())

	def recv(self, buf):
		return self.sokcet.recv(buf)

	def card_connected(self, tag):
		if tag.type != "Type3Tag":
			return False

		try:
			tag.idm, tag.pmm = tag.polling(self.system_code)
		except nfc.tag.tt3.Type3TagCommandError as err:
			log.error("polling error: " + err)
			return False

		for key in sorted(self.service_code.keys()):

			service_code = self.service_code[key]
			sc_list = [nfc.tag.tt3.ServiceCode(service_code >> 6, service_code & 0x3f)]
			try:
				if tag.request_service(sc_list) == [0xFFFF]:
					continue
			except nfc.tag.tt3.Type3TagCommandError as err:
				continue
			
			if isinstance(tag, nfc.tag.tt3_sony.FelicaStandard) or isinstance(tag, nfc.tag.tt3_sony.FelicaMobile):
				try:
					bc = nfc.tag.tt3.BlockCode(0x00, service=0)
					data = tag.read_without_encryption(sc_list, [bc])
					log.info("block: " + binascii.hexlify(data))
				except Exception as e:
					log.error("card read error: " + e)

				if key == "suica":
					print("suica balance: (little endian)%s" % binascii.hexlify(data[10:12]))
				elif key == "univ":
					self.send(binascii.hexlify(data[0:6]))
				elif key == "edy":
					self.send(binascii.hexlify(data[2:10]))
				elif key == "waon":
					# self.send(binascii.hexlify(data[2:10]))
					print("waon balance: %s" % binascii.hexlify(data[0:24]))
				elif key == "nanaco":
					self.send(binascii.hexlify(data[0:8]))
				else: 
					log.error("error: tag isn't Type3Tag")
				break

		return True

	def open(self, e):
		topen = threading.Thread(name='open', target=self.turn_on, args=(e,))
		tclose = threading.Thread(name='close', target=self.turn_off, args=(e,))

		topen.start()
		time.sleep(2)
		e.set()
		
		tclose.start()
		time.sleep(2)
		e.set()

	def turn_on(self, e):
		while not e.isSet():
			self.piface.output_pins[1].turn_on()
		e.clear()
		
	def turn_off(self, e):
	        while not e.isSet():
	                self.piface.output_pins[1].turn_off()
	        e.clear()

 
	def run(self):
		try:
			clf = nfc.ContactlessFrontend('usb')
		except IOError as error:
			raise SystemExit(1)
		try:
			return clf.connect(rdwr={'on-connect': self.card_connected})
		finally:
			clf.close()


if __name__ == '__main__':
	e = threading.Event()
	reader = Reader('/tmp/laputa.sock')
	reader.socket_connect()
	
	while reader.run():
		# True: "I can't see!! I can't see!!!!"
		# False: "Get down on your knee. Beg your life."
		data = reader.recv(40)
		if data == "I can't see!! I can't see!!!!":
			reader.open(e)

		log.info("*** RESTART ***")
