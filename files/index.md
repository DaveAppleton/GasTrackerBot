
🎩 Gas Tracker Bot {{.Version}} 🎩
``` 
         ETHGasStation
--------+--------------
Fastest | {{.EGS.Fastest.StringFixed 2}}
Medium  | {{.EGS.Fast.StringFixed 2}}
SafeLow | {{.EGS.SafeLow.StringFixed 2 }}
--------+--------------
         Sparkpool
--------+--------------
Fastest | {{.GN.Fast.StringFixed 2}}
Medium  | {{.GN.Medium.StringFixed 2}}
SafeLow | {{.GN.Safe.StringFixed 2}}
--------+--------------
         Etherscan
--------+--------------
Fastest | {{.ES.FastGasPrice.StringFixed 2}}
Medium  | {{.ES.ProposeGasPrice.StringFixed 2}}
SafeLow | {{.ES.SafeGasPrice.StringFixed 2}}

```

{{.Advert}}