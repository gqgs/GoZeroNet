# GoZeroNet
Go implementation of the [ZeroNet](http://zeronet.io/) network.

This aims to be compatible with the the _core_ of the [original](https://github.com/HelloZeroNet/ZeroNet) Python implementation.


### Status

- [Network protocol](https://zeronet.io/docs/help_zeronet/network_protocol/)
    - [x] handshake
    - [x] getFile
    - [x] streamFile
    - [x] ping
    - [x] pex
    - [x] update
    - [x] listModified
    - [ ] getHashfield
    - [ ] setHashfield
    - [x] findHashIds
    - [x] checkport
    - [ ] getPieceFields
    - [ ] setPieceFields
- [ZeroFrame API](https://zeronet.io/docs/site_development/zeroframe_api_reference/)
    - [x] Wrapper commands
    - [ ] as
    - [x] announcerInfo
    - [x] announcerStats
    - [x] certAdd
    - [x] certSelect
    - [x] channelJoin
    - [x] channelJoinAllsite
    - [x] dbQuery
    - [ ] dirList
    - [x] fileGet
    - [x] fileList
    - [x] fileNeed
    - [ ] fileQuery
    - [ ] fileRules
    - [x] fileWrite
    - [x] ping
    - [x] siteInfo
    - [x] siteList
    - [ ] sitePublish
    - [ ] siteReload
    - [ ] siteSign
    - [x] siteUpdate
    - [x] siteSetLimit
    - [x] userGetSettings
    - [x] userSetSettings
    - [x] userGetGlobalSettings
    - [x] serverShutdown
- Announcers
    - [X] HTTP
    - [x] PEX
    - [ ] UDP
    - [ ] Zero
- Files
    - [x] Required
    - [x] Optional
    - [x] Big
- [ ] UPnP
- [ ] Hidden service
- [ ] TLS
- [ ] Namecoin


### Is it safe already to use for unstrusted sites?
Not yet!
