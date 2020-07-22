import sys
import glob


sys.path.append('gen-py')
sys.path.append('gen-py/commu')


from commu import ttypes
from commu import Calculator
from commu.ttypes import InvalidOperation, Operation, Work


import thrift
from thrift import Thrift
from thrift.transport import TSocket
from thrift.transport import TTransport
from thrift.protocol import TBinaryProtocol
from thrift.protocol import TJSONProtocol

def main():
    # Make socket
    transport = TSocket.TSocket('192.168.20.45', 8091)

    # Buffering is critical. Raw sockets are very slow
    transport = TTransport.TBufferedTransport(transport)

    # Wrap in a protocol
    # protocol = TBinaryProtocol.TBinaryProtocol(transport)
    protocol = TJSONProtocol.TJSONProtocol(transport)
    # Create a client to use the protocol encoder
    client = Calculator.Client(protocol)

    # Connect!
    transport.open()

    client.ping()
    print('ping()')

    sum_ = client.add(1, 1)
    print('1+1=%d' % sum_)

    work = Work()

    work.op = Operation.DIVIDE
    work.num1 = 1
    work.num2 = 0

    try:
        quotient = client.calculate(1, work)
        print('Whoa? You know how to divide by zero?')
        print('FYI the answer is %d' % quotient)
    except InvalidOperation as e:
        print('InvalidOperation: %r' % e)

    work.op = Operation.SUBTRACT
    work.num1 = 15
    work.num2 = 10

    diff = client.calculate(1, work)
    print('15-10=%d' % diff)

    log = client.getStruct(1)
    print('Check log: %s' % log.value)


    a = 0
    if sys.argv[1]=='1':
        a = client.doconfig(1, 1, 'start filebeat')
    elif sys.argv[1] == '2':
        a = client.doconfig(2, 2, 'stop filebeat')
    elif sys.argv[1] == '3':
        a = client.doconfig(3, 3, 'start packetbeat')
    elif sys.argv[1] == '4':
        a = client.doconfig(4, 4, 'stop packetbeat')
    elif sys.argv[1] == '5':
        a = client.doconfig(5, 5, '{"filebeat":{"input":"myinput"}}')



    # Close!
    transport.close()


if __name__ == '__main__':
    try:
        main()
    except Thrift.TException as tx:
        print('%s' % tx.message)




