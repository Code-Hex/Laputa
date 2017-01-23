# Laputa
Key certification system for Department of Information Engineering, University of the Ryukyu.

```
                      （    ヽ                                     ,  ⌒ヽ  , ⌒ヽ
         ,  ⌒ヽ    (          ヽ     _,=''''''^~~~~~~~~~^''''=,,,,（        （        '
       （        '  （            ,-='''~  -=^~~~^-^~~~^==-  '=,,,        ゝ        ｀ヽ
        ゝ        ｀ヽ        <~  -==^~~~^ =^~~^=-=^~~^'=-~'=,    (                ）
      (                        ヽ'^' __,,,,,,i~~~l===|~~i==|~~|＿,,,,,..ノ  （                    ヽ
    (        (⌒                ヽi~ |   |＿_レ、l--l--レ.;---i  i-、  （                      ）
  （                            r'^~~~~l l     | :| ∩ ∩|,-=,__,-,_|  |~i^i,,                      ヽ
    （                ｀)     l^^|,,,,--==.i~~l~~~~~~~~~~|    i  l .|~^''''l~^i,,,,                      ヽ
      ,ゝ                 ／i~~i' l ∩∩l  .ｌ  ∩  ∩   l    |__| .| .∩|  .|  l-,
    （        '      ,,,,,='~|  |  |' |,,=i~~i==========|~~|^^|~ ~'i----i==i,, | 'i
  （       （        |  l ,==,-'''^^   l   |. ∩. ∩. ∩. |   |∩|    |∩∩|   |~~^i~'i、
（               ,=i^~~.|   |.∩.∩ |,...,|__|,,|__|,,|__|,,|__|,....,||,,|.|,.....,||,|_|,|.|,....,|    |  |~i
  （            l~| .|    |  ,,,---== ヽノ       i        ヽノ~~~ ヽノ     ~ ソ^=-.i,,,,|,,,|
    ,ゝ       .|..l  i,-=''~~--,,,    ＼   ＼   l      ／      ／        ／    __,-=^~
  （          |,-''~  -,,,＿    ~-,,.   ＼ .＼ |  .／     ／    ＿,,,-~     ／          ヽ
,  ⌒ヽ       ~^''=、＿ _ ^'- i=''''''^~~~~~~~~~~~~~~~~~~~~^''''''''=i -'^~  （            ヽ
         ヽ                ~^^''ヽ  ヽ   i     |    ｌ   i   /   ／   ノ      （                 ）
      ,  ⌒ヽ                    ヽ    、 l   |   ｌ   l  / .／    ／     （                  （
          ,  ⌒ヽ⌒ヽ             ＼_  、i ヽ   i   /      ,,=='      （                      ヽ
                    ,  ⌒ヽ            ''==,,,,＿＿＿,,,=='~            ,ゝ                      ヽ
          ⌒                                  （                            ヽ
      ノ                                       ヽ
```

# What is this?
It is a felica reader that exists in the 1st building of the Faculty of Engineering Building of the University of the Ryukyus.

# Logic
```
                              check
                            felica id (felica id searcher)
           +--------------+            +---------------+
           |              <--------+   |               |     felica data
           |   Akatsuki   |        |   |   Levistone   <------------------+
           |              +----+   |   |               |                  |
           +--+---+----+--+    |   |   +-------+-------+                  |
              |   ^    |       |   |           |                          |
              |   |    |       |   |           |                 +--------+-------+
initialize    |   |    |       |   |           |                 |                |
    to        |   |    |       |   |           |                 |   nfc reader   |
authenticate  |   |    |       |   |           |                 |                |
              |   |    |       |   |           | send felica id  +----------------+
              |   |    |       |   |           |
           +--v---+----v--+    |   |   +-------v-------+
           |              |    |   +---+               |
           |    Laputa    |    |       |     Balus     +-------------------------->
           |              |    +------->               |
           +--------------+            +---------------+      open the door
                              recieve
         (RESTful API Server)         (Unix domain socket)
                   +                            ^
                   |                            |
                   |                            |
                   |      +-------------+       |
                   |      |             |       |
                   +------>   leveldb   +-------+
             put secret   |             |   get secret
                          +-------------+
```

# Build
There are build modes for development and staging.
- for develop
    
        make build-dev

- for staging

        make build-staging

# Setup
After `git clone` this project

    ./migrate

# Run
After setup, you can run

    carton exec ./run

# Restart
You can run without finish after program is modified(Golang projects only).

    make restart-staging

# Configuration
Please open `run` with your favorite editor.  
By editing the env function, you can change the behavior at program execution.  
  
Environment variable
- `LAPUTA_CERTFILE` Specify the file path of the certificate
- `LAPUTA_KEYFILE` Specify the file path of the key
- `LAPUTA_AKATSUKI` Specify the URL of the api for authentication
- `LAPUTA_FLOOR` Floor code for registration
