import subprocess
from subprocess import call
from subprocess import Popen,PIPE
from time import sleep
import os


ifconfig = subprocess.check_output('ifconfig')
ifconfig = ifconfig.split(b"\n\n")
for i in ifconfig:
    if(i.startswith(b'tun')):
        tun = i.split(b":")[0]
        subprocess.run(["sudo","ifconfig",tun,"down"])

if os.path.exists("vpn_connection.txt"):
        os.remove("vpn_connection.txt")
p = subprocess.Popen('gnome-terminal -- bash -c "sudo openvpn starting_point_vpn.ovpn > vpn_connection.txt; exec bash"',shell=True, stdout=subprocess.PIPE, universal_newlines=True)

file_present = False
while file_present == False:
    if os.path.isfile('vpn_connection.txt'):
        file_present = True
        break
    sleep(1)

    
isc = False
while(isc == False):
    with open('vpn_connection.txt') as f:
        lines = f.readlines()
        if "Initialization Sequence Completed" in lines[-1]: 
            print(lines[-1])
            isc = True
            break
    sleep(1)

ifconfig = subprocess.check_output('ifconfig')
ifconfig = ifconfig.split(b"\n\n")
for i in ifconfig:
    if(i.startswith(b'tun')):
        tun = i.split(b":")[0]
        print(tun)
