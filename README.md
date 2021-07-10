# GoZeroNet
Go implementation of the [ZeroNet](http://zeronet.io/) network.

This aims to be as compatible with the the _core_ of the [original](https://github.com/HelloZeroNet/ZeroNet) Python implementation with the exception of some minor changes. These changes are listed bellow.

- `dbQuery` returns an error if it tries to select multiple columns with the same name. This behavior can be disabled by setting `validate_database_queries` to false in `zeronet.toml`.
- The regular expressions defined in content.json files must be accepted by the [parser used in the Go standard library](https://golang.org/pkg/regexp/). That means a number of features are not supported to guarantee a linear time complexity. For more details read the [RE2 syntax definition](https://github.com/google/re2/wiki/Syntax).

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
    - [x] fileDelete
    - [x] fileGet
    - [x] fileList
    - [x] fileNeed
    - [ ] fileQuery
    - [ ] fileRules
    - [x] fileWrite
    - [x] ping
    - [x] siteInfo
    - [x] siteList
    - [x] sitePublish
    - [ ] siteReload
    - [x] siteSign
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
    - [x] Big (download)
    - [x] Big (upload)
- DBSchema mappings
    - [x] to_table (string array)
    - [ ] to_table (object array)
    - [x] to_json_table
    - [ ] to_keyvalue
- [ ] ~~UPnP~~
- [ ] Hidden service
- [ ] TLS
- [ ] Namecoin


### Is it safe already to use for unstrusted sites?
Not yet!
