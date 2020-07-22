
from pysnmp.hlapi import *

g= sendNotification(SnmpEngine(),CommunityData('public', mpModel=0),UdpTransportTarget(('127.0.0.1', 9000)), ContextData(), 'trap',        NotificationType( ObjectIdentity('1.3.6.1.4.1.20408.4.1.1.2.0.432'),).addVarBinds(('1.3.6.1.2.1.1.1.0', OctetString('my system')) ))

next(g)

