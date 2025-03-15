# Bitcoin
Bitcoin Core Implementation in Golang

## Elliptic Curve
Bitcoin relies heavily on a match object called elliptic curve,
without this math structure bitcoin will like a castle on a beach, 
it will collapse at any time.
What is an elliptic curve, it's an equation like this: y^2 = x^3 + ax +b, and its shape just like the following:

- For bitcoin, its elliptic curve has a name: SECP256K1 and its equation is y^2 = x ^ 3 + 7.

## Signature and Verify
For signature, we can't use SEC for encoding because there is no close relation between r and s,
you know one of them, but you can't derive the other.
There is a scheme to encode signature
called DER (Distinguished Encoding Rules), following are steps for encoding signature:

1. set the first byte to 0x30

2. the second byte is the length of signature (usually 0x44 or 0x45)

3. the third byte set to 0x02, this is an indicator that the following bytes are for r.

4. transfer r into bytes array, if the first byte of r is >= 0x80, then we append a byte 0x00 in the beginning of the bytes, compute the length of the bytes array, append the value of length
after the marker byte 0x02 in step 3, and append the bytes array following the value of length

5. add a marker byte 0x02 at the end of byte arrays from step 4.

6. append s as how we append r in step 4.

We need to encode the length of r and s, because r and s at most has 32 bytes,
but sometimes their length may be shorter than this.
Let's see an example for the encoding of r and s:

30 45 02 21 00 ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f 02 20 7a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed

the first byte is 0x30 as mentioned above, the second byte is 0x45 that is the total length of r and s. The third byte is marker 0x02 indicating from next byte is the beginning of bytes array for r.
According to 3, the byte following from marker 0x02 is the length of bytes array for r, the value is 0x21. Byte following 0x21 is 0x00, by step 4, it indicates the first byte of r is >= 0x80, we can
see the byte follow 0x00 is 0xed which is indeed more than 0x80, the length of r is 0x21 - 1 = 0x20, which means the following 32 bytes are bytes array for r, we extract it out as following:

r: 81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f

following the last byte of r is 0x20 , it is indicator for the beginning of s, the byte following 0x02 is 0x20 which indicates the length of s is 0x20, the following byte is not 0x00, which means the
first byte of s is not more than 0x80, and the byte following 0x20 is 0x7a , this is the beginning byte for s and is smaller than 0x80, the byte array for s is :

s: 7a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed

## Wallet Address
As a bitcoin user, we always need to send or receive bitcoins from others;
this will require letting others know your wallet address. 
Because wallet address needs to be read by human, all the encoding schemas 
we have before are produce result in binary,
therefore we need another scheme to create wallet address in a human friendly way.

Wallet address is actually generated from a public key, and it needs to satisfy the following requirement:

1. Easy to read and write, user may want to memorize it or write it down on paper

2. Not too long for sending over the internet

3. It should be secure, and harder to make mistake, you don't want you fund transfer to people unknown to you!

The base58 encoding scheme can help us to achieve three goals.
Compared with the common use of base64, it removes characters like l,
I, 0, O, -, _ because they are easy to confuse with each other.
Because the encoding schema uses all numbers, and uppercase and lowercase letters and remove 0 O,
l, I, which means it will use 58 characters in the encoding process.


## Transaction Validation
For a bitcoin node, one of its major tasks is to validate a transaction;
there are several steps to take for it, the first thing is to check the output can match to the transaction.
For example, if a transaction
is about "Jim using 10 dollars to buy a cup of coffee with price of 3 dollars," then we need to check :

1. Jim really has 10 dollars

2. the amount left after buying the coffee should be 7 dollars

If the transaction is honest, then the input of the transaction(10 dollars) should be greater than the output of the transaction(7 dollars), that is when we use the amount of input minus the amount of the output
the result should be positive, if the result is negative, then the transaction is "dishonest" it wants to fake money from air.

## Transaction Creation
A transaction is recording an event of bitcoin transition,
it needs to make sure where the bitcoins transfer to, is the input for this transaction legal or valid,
and how quickly the transaction can on chain which
means the transaction is legally accepted.

Let's see how to construct a valid transaction and send it to the network,
this process maybe like you first go to the bank,
deposit some amount of money in your account and transfer some of them to your friend.
The first thing we need to do is convert a wallet address from base58 encoded,
that is we need a process that can get its original content when the input is encoded by base58.

## Blocks of Transactions for BitCoin Chain Block

In the previous section, we successfully broadcast our transaction on to the bitcoin blockchain.
Since at any given time, there are thousands of transactions waiting for broadcast to chain,
every bitcoin node will batch all coming transactions at
broadcast them at every ten minutes.
The collection of transactions that are broadcast at one time is called block.


We first pay attention to a special transaction that is the first transaction in the block,
and it has the name of coinbase transaction.
There is a star bitcoin company called coinbase and is already listed on Nasdaq, but the coinbase transaction
we are going to look has nothing to do with it.
The coinbase transaction is every important for bitcoin nodes,
because nodes can get significant reward.Let's see an example of coinbase transaction,
following is the raw binary data of one coinbase
transaction:
```g
01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff5e03d71b07254d696e656420627920416e74506f6f6c20626a31312f4542312f4144362f43205914293101fabe6d6d678e2c8c34afc36896e7d9402824ed38e856676ee94bfdb0c6c4bcd8b2e5666a0400000000000000c7270000a5e00e00ffffffff01faf20b58000000001976a914338c84849423992471bffb1a54a8d9b1d69dc28a88ac00000000
```
Let's dissect the data block above piece by piece:
1. the first four bytes: ``01000000`` it is the version of the transaction in little endian format

2. the following one byte: 01 is the input count

3. the following chunk of zeros: ``0000000000000000000000000000000000000000000000000000000000000000``,
   is the previous transaction hash,
since its first transaction of the block, therefore, it has not previous transaction, and this value is all 0,

4. ``ffffffff`` previous transaction index

5. the following data chunk:
`` 5e03d71b07254d696e656420627920416e74506f6f6c20626a31312f4542312f4144362f43205914293101fabe6d6d678e2c8c34afc36896e7d9402824ed38e856676ee94bfdb0c6c4bcd8b2e5666a0400000000000000c7270000a5e00e00``
is input script

6. ``ffffffff`` sequence number

7. 01 output count

8. ``faf20b5800000000`` output amount

9. `` 1976a914338c84849423992471bffb1a54a8d9b1d69dc28a88ac`` p2pkh scriptPubKey

10. ``00000000`` lock time

The structure of coinbase transaction is the same as we have seen before, but has some specials:
1. coinbase transaction must have exactly one input
2. the one input can only have previous transaction id set to 32 bytes of data chunk and filled with all 0
3. the one input can only have previous transaction output index of four bytes with each byte set to value 0xff

## Block Header

Block in blockchain is like packet for the internet,
the raw data of transaction is like payload of a packet for the internet.
Since an internet packet has a header, and block also has its own
header.
Let's dissect block binary data into fields, following is an example of raw data for a block:

```ggg
020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d
```

1. the first four bytes is version in little endian, 02000020

2. following 32 bytes of data chunk in little endian indicate the previous block:
   ``8ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd000000000000000000``,
   this is the hash254 result of the previous block.

3. the following 32 bytes in little endian is the make root,
``5b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be``,
we will go to this concept later on.

4. the following four bytes is timestamp in little endian format: 1e77a75

5. the following four bytes are called bits: e93c0118

6. the final four bytes are called nonce: a4ffd71d

Compare with transaction, filed in block header are in fix length, and the total length of block header is 80 bytes.

## Proof of Work (Blocks of Transaction)

When miners broadcast a block on to the chain,it needs to pay some cost for it,
the cost is to prevent adverse node trying to create fake
transactions and group them in a block and add the chain.
This is like when a gangster wants to accept a member, the boss will ask the guy
to commit some crime like killing somebody; by doing this can prevent the guy beginning an undercover of the police.

the miners need to do some heavy computation, when it has the result, other miners can easily verify the result,
and then the given miner will allow putting its block on to the chain.

Let's see the details of proof-of-work, in a previous section,
we compute the hash256 of the given block header and have the following result:
```ggg
0000000000000000007e9e4c586439b0cdbe13b1370bdd9435d76a644d047523
```


We can see there are several zeros at the beginning of the hash result.
Since the result of hash256 is very random, therefore given any bit in the result, the probability that it takes
the value of zero is 0.5, a zero in the above result needs four bits to represent, 
which means a 0 happening at the result has probability of 0.5*0.5*0.5*0.5.
The probability of more zeros appearing in the result will decrease dramatically.
As there are 18 zeros appear at the beginning of the result,
then its probability is (0.5)^(4*18) which is roughly 10^22.


If one wants to get the same result above, it needs to do 10^22 rounds of computation, which is a huge computation task.
There is a concept
called difficulty in bitcoin mining,
which is how many zeros need to appear at the beginning of block header hash, when the difficulty
increases, the more zeros need to appear at the beginning of the hash.


The question is how to make a given number of zeros appear at the beginning of hash?
This is where the nonce comes in.
The miner can manipulate
this field to generate the hash to meet the requirement.
Since there are 4 bytes for the nonce field,
which means the miner can try 2^32 times to get the required hash.


The number of zeros that are required at the beginning of the hash is determined by the value
computed from the bit field.
The computed
value has the name of "target," as long as the hash result is smaller than the target;
then it will guarantee there are enough zeros at the
beginning of the hash256 result.


In order to compute the target, we need to separate the bits filed into two parts, since it has 4 bytes, 
then we separate it to two parts by
using the last byte as the first part which is called exponent, 
and the last three bytes is another part call coefficient, then the value of
target is :

```ggg
target = coefficient * 256^(exponent - 3)
```

We need to get more explanation for the above equation.
The bit field in the block header is used to encode a kind of hex value that are only
three no zero digits at the left and lots of zero on the right, such as the following value:

```ggg
0000000000000000013ce9000000000000000000000000000000000000000000
```

We can see that there are only three no zero bytes which are 0x01, 0x3c, 0xe9,
and there are dozens of zeros on the right.
We can encode
such value by using four bytes, we can use three bytes to encode the three no zeros bytes,
and the last byte can be used to record the total
number of bytes in the value.

Notices that two digits are corresponding to one byte,
therefore the value of ``0x013ce9000000000000000000000000000000000000000000`` has 0x18 which
is 24 bytes in total, then we can encode the above value in four bytes as following:

```ggg
18 01 3c e9
```

if we want to get the value from above encoding, we can first take the last three bytes and assembly them into hex value of 0x013ce9,
then we need to left shift this value with 0x15 bytes, each byte has 8 bits, which means we need to do is:
0x013ce9 << 8 * 0x15

since a value shift left one bit is like the value multiply the value of 2,
therefore the equation above is equivalent to:

```ggg
0x013ce9 * 2^(8*15) => 0x013c39 * (2^8)^(15) => 0x013ce9 * 256^(0x15) => 0x013ce9 * 256^(0x18-3)
```

And notice that if we reverse the byte order of 18 01 3c e9,
then we get e9 3c 01 18 this is exactly the value of bits field in our
block header.
