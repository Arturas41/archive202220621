package windows

var driverVioGPUDo = DriverInfo{
	PackageName: "viogpudo.inf_amd64_8224060246e67964",
	DriversRegistry: `[\DriverDatabase\DeviceIds\{{ classGuid|lower }}]
"{{ infFile }}"=hex(0):

[\DriverDatabase\DeviceIds\pci\VEN_1AF4&DEV_1050&SUBSYS_11001AF4&REV_01]
"{{ infFile }}"=hex(3):01,f9,00,00

[\DriverDatabase\DriverInfFiles\{{ infFile }}]
@=hex(7):{{ packageName|toHex }},00,00,00,00
"Active"=hex(1):{{ packageName|toHex }},00,00

[\DriverDatabase\DriverPackages\{{ packageName }}]
@=hex(1):6f,00,65,00,6d,00,30,00,2e,00,69,00,6e,00,66,00,00,00
"Catalog"=hex(1):76,00,69,00,6f,00,67,00,70,00,75,00,64,00,6f,00,2e,00,63,00,61,00,74,00,00,00
"ImportDate"=hex(3):80,13,e8,66,f6,9d,d7,01
"InfName"=hex(1):76,00,69,00,6f,00,67,00,70,00,75,00,64,00,6f,00,2e,00,69,00,6e,00,66,00,00,00
"OemPath"=hex(1):5c,00,5c,00,31,00,39,00,32,00,2e,00,31,00,36,00,38,00,2e,00,31,00,37,00,38,00,2e,00,37,00,30,00,5c,00,73,00,68,00,61,00,72,00,65,00,64,00,5c,00,76,00,69,00,72,00,74,00,69,00,6f,00,5c,00,76,00,69,00,6f,00,67,00,70,00,75,00,64,00,6f,00,5c,00,77,00,31,00,30,00,5c,00,61,00,6d,00,64,00,36,00,34,00,00,00
"Provider"=hex(1):52,00,65,00,64,00,20,00,48,00,61,00,74,00,2c,00,20,00,49,00,6e,00,63,00,2e,00,00,00
"SignerName"=hex(1):4d,00,69,00,63,00,72,00,6f,00,73,00,6f,00,66,00,74,00,20,00,57,00,69,00,6e,00,64,00,6f,00,77,00,73,00,20,00,48,00,61,00,72,00,64,00,77,00,61,00,72,00,65,00,20,00,43,00,6f,00,6d,00,70,00,61,00,74,00,69,00,62,00,69,00,6c,00,69,00,74,00,79,00,20,00,50,00,75,00,62,00,6c,00,69,00,73,00,68,00,65,00,72,00,00,00
"SignerScore"=dword:0d000005
"StatusFlags"=dword:00000512
"Version"=hex(3):00,ff,09,00,00,00,00,00,68,e9,36,4d,25,e3,ce,11,bf,c1,08,00,2b,e1,03,18,00,40,ef,05,7a,77,d7,01,b0,4f,68,00,55,00,64,00,00,00,00,00,00,00,00,00

[\DriverDatabase\DriverPackages\{{ packageName }}\Properties]

[\DriverDatabase\DriverPackages\{{ packageName }}\Properties\{4da162c1-5eb1-4140-a444-5064c9814e76}]

[\DriverDatabase\DriverPackages\{{ packageName }}\Properties\{4da162c1-5eb1-4140-a444-5064c9814e76}\0009]
@=hex(ffff0012):33,00,30,00,30,00,39,00,37,00,37,00,37,00,30,00,5f,00,31,00,34,00,31,00,35,00,35,00,36,00,33,00,31,00,34,00,35,00,36,00,37,00,30,00,36,00,35,00,39,00,32,00,5f,00,31,00,31,00,35,00,32,00,39,00,32,00,31,00,35,00,30,00,35,00,36,00,39,00,33,00,36,00,38,00,33,00,38,00,34,00,38,00,00,00
`,
}
