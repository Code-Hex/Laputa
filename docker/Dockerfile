FROM raspbian
MAINTAINER Code-Hex <x00.x7f@gmail.com>

RUN apt-get update && apt-get install -y apt-utils sudo && apt-get -y dist-upgrade
RUN apt-get install -y python python-pip python3 python3-dev python3-pip

RUN mkdir /var/run/sshd

RUN useradd pi \
    && mkdir /home/pi \
    && chown pi:pi /home/pi \
    && addgroup pi staff \
    && addgroup pi sudo \
    && true

RUN chsh -s /bin/bash pi
RUN echo 'pi:raspberry' | chpasswd
RUN mkdir -p /home/pi/.ssh && chmod go-rwx /home/pi/.ssh

# FileCopy(local to container image)
COPY ./authorized_keys /home/pi/.ssh/authorized_keys
COPY ./bashrc /home/pi/.bashrc

RUN chmod go-rw /home/pi/.ssh/authorized_keys && chown -R pi:pi /home/pi

# Enable passwordless sudo for users under the "sudo" group
RUN sed -i.bkp -e \
      's/%sudo\s\+ALL=(ALL\(:ALL\)\?)\s\+ALL/%sudo ALL=NOPASSWD:ALL/g' \
      /etc/sudoers

USER pi
WORKDIR /home/pi

EXPOSE 22
CMD ["sudo", "/usr/sbin/sshd", "-D"]
