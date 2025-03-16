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
Since there are four bytes for the nonce field,
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

## P2P Networking
One of the great creativity of bitcoin blockchain is it is a distributing system.
Thousands of independent nodes can work together just like an
integrated system, and even different nodes may be far away from each other, 
the system can still guarantee no any one can sabotage the whole
system, and the nodes in the system can make sure their data can synchronize with each other and make sure of the data integrity.


All these achievements are thanks to the bitcoin network protocol,
we will dive deep into the bitcoin networking protocol and make clear how
such magic is happening.
Following is an example of bitcoin networking package raw data:

```ggg
f9beb4d976657273696f6e0000000000650000005f1a69d2721101000100000000000000bc8f5e5400000000010000000000000000000000000000000000ffffc61b6409208d010000000000000000000000000000000000ffffcb0071c0208d128035cbc97953f80f2f5361746f7368693a302e392e332fcf05050001
```
Let's dissect it into fields:

1. the first 4 bytes are always the same and referred to as network magic number: ``f9beb4d9``,
   its usage is to tell receiver that, when you see
these four bytes appear together, then you should know this is the beginning for bitcoin networking package.
   And this number used to
differentiate the main-net with testnet, for testnet, the four bytes are 0b110907.

2. the following twelve bytes are called command: ``76657273696f6e0000000000``,
   it used to describe the purpose of this packet.
   It can be human-readable string.

3. the following four bytes: ``65000000`` it is a payload length in little endian format,

4. the following four bytes: ``5f1a69d2`` is the first 4 bytes of hash256 of the payload.

5. the following bytes are data of the payload


We can use the channel to send more requests to
the peer to ask more data.
The most request data is requesting block header info from a full node.
When you have block headers, then you can
request the block body from multiple peers or do many useful works by those headers.


As we have seen the version command before,
the command for getting block headers has name ``getheaders`` and following is raw binary data for
payload of ``getheaders`` command:

```ggg
7f11010001a35bd0ca2f4a88c4eda6d213e2378a5758dfcd6af437120000000000000000000000000000000000000000000000000000000000000000000000000000000000
```

Let's dissect the above data chunk into fields as following:

1. the first four bytes is protocol version as we have seen before: 7f110100

2. the following bytes are variant used to indicate the amount of hash, here the value is 01 which is only one byte

3. the following bytes are the starting block hash we want to request:

``a35bd0ca2f4a88c4eda6d213e2378a5758dfcd6af43712000000000000000000``

4. the following bytes are the ending block hash we want to request:

``0000000000000000000000000000000000000000000000000000000000000000``

All zeros here mean we want to get as block headers as many as possible, but the maximum number we can get is 2000,
which is almost a difficulty adjustment period (2016 blocks).


## Merkel Tree
All kinds of blockchain suffer a common weakness that is shortage of disk volume, computing power and bandwidth,
bitcoin blockchain is not an exception.
This causes problems when you want to check the data
integrity, you may need to wait for a long time to download the necessary data for a simple checking.
For example, when we're using bitcoin for trading, you pay your goods or services with bitcoins, the
paying record needs to append to the bitcoin chain.


In order to verify that your transaction is recorded on the chain,
you may need to sync the whole chain down to your device,
this may need weeks and of course unacceptable.

Therefore, we need some algorithm
to quickly verify that the given info is already on the chain.
To make sure some data in a block and exist on the chain is called proof-of-inclusion.


Merkle tree is a kind of data structure used to verify the integrity of a group of objects,
in blockchain it is used to verify transactions
in block, the objects in the group are ordered, then we can use the following steps to construct merkle tree:

1. Hash each object in the group using given hash function.

2. If there is only one object to be hashed, then it is after the hash of the process is completed.

3. if there is an odd number of objects, then after hashing them all,
we copy the last hash result and add it to the hashed list
then we have even numbed of hash values in the list.

4. we select a pair of hash results in order,
hash the pair and use the result as their parent, this step will half the number of items in the second layer.

5. goto step 2

This is somehow like a reverse of a binary tree, we begin from the lowest level,
then select two leaves and construct their parent until
we go to the root.


The process is just like divide and conquer, if there are too many objects, hash them together is impossible, 
then we can divide the group into two, 
and try to hash each half group then combine the two hash results to produce the top hash, 
if the subgroup is still too large,
then we continue to divide them until we get only one object, 
then we collect the lower level hash to produce the up lever hash.


We can abstract the process as following:
H: hash function
P: parent hash
L: left hash
R: right hash

Then ``P = H(L || R)``, (|| means connecting R to the end of L).

If someone wants to proof L is included in P, then they can provide R and P,
then we compute L, then connect L and R to compute the hash and check the result is P or not.


Using merkle tree, we can achieve proof of inclusion, for example, you have a list of objects , 
and you want to check your peer has the same set of objects, then you compute the merkle root on L and
let the peer compute the merkle root, then you compare two roots, if they are the same, 
then it is proof your peer also have the same list of objects.
In application to the bitcoin, if we want to know
whether a batch of transactions has included in the block or not,
we can compute the merkle root instead of getting all transactions from the block and checking them one by one.


If we have two transactions ``H(k), H(N)``,
and we want to make sure whether these two transactions have already been included in a block, 
then the bitcoin full node can only return hashes that are represented by the
blue box.
We can compute the merkle root for a check in the following step:

1. ``H(KL) = MerkleParent(H(k), H(L))``

2. ``H(MN) = MerkleParent(H(M), H(N))``

3. ``H(IJKL) = MerkleParent(H(IJ), H(KL))``

4. ``H(MNOP) = MerkleParent(H(MN), H(OP))``

5. ``H(IJKLMNOP) = MerkleParent(H(IJKL), H(MNOP))``

6. ``H(ABCDEFGHIJKLMNOP) = MerkleParent(H(ABCDEFGH), H(IJKLMNOP))``


As we can see, the merkle root is a kind of compression algorithm;
we need only part of the information can we get to the conclusion.
But we still have problems with the above computation, How do we know we
need to pair H(k), H(L) to get the merkle parent, how do we know we need to pair H(NM), H(OP) to get merkle parent?
We need some info th deduce such info.


The info we need is the position of those blue boxes,
and we need to define the "position" of the box in the binary tree structure.
The "position" of a node in a tree related to how
we "travel" the tree, if you have background of basic data structure and algorithm, 
you will know there is a kind of data structure called "graph", the tree above is binary tree, and it is a kind of graph.


For graph which contains "nodes and edges," if given one node,
we can go to other nodes by using the edges coming from the given node,
then we have two kinds of ways to "travel" the graph, breathe first
and depth first.


As you can see from the image above, for breath first travel, we visit the node layer by layer, we first at the root node 0, then we visit all nodes below it that are node 1, 2, 3, then we goto layer 2,
visit all nodes there that are nodes 4,5,6,7.
The order for each node that they are visited is their "position."
The depth-first travel will be a little bit complex, when we in a given node, we first check
if we can go to the lower layer from the left edge, if we can, then we go down to the lower layer.
If we can't go down to the lowe layer from the left edge(there is not left edge, or we have already been there),
then we try whether we can go down by using the right edge, if we can then we go down to the lowe layer by using the right edge.
Otherwise, we go up to the parent node and do the same again, the order for
the node that is being visited is the "position" of that node,
you can see nodes will have different order or "position" for these two travel ways.


When the full node returns hashes for the blue boxes, it will also return their order under the depth-first travel,
then we will use the order and the given hash value to reconstruct the merkle root.
Let's go through the whole process step by step:


the first step we need to achieve is given a list of objects, we need to convert it to tree like structure.


As you can see from the above image, we have 8 nodes in the list. 
if we want to construct a merkle tree from these 8 nodes, 
we can set these 8 nodes as the leaves, and pair two as a group then "grow"
its parent, therefore for 8 nodes, we can "grow out" 4 parents in the second layer, 
the same process can repeat again and again until we have only one node. 
Since the number of nodes is half when it comes up to one layer, given N nodes, we can have most int(lg(N)+1) layers.


## Merkel Block

We have built up merkle tree in the previous section,
in this section we will see how to use the tree to verify proof-of-inclusion in bitcoin blockchain.
As we have shown in the previous section,
When we get a list of hashes, we can build up the merkle tree and get the root value.
The question is how we get those lists of hash values, actually we can send a ``getdata`` command with
given transaction hash to a full node, then it will respond with a ``merkleblock`` command,
and a list of hash values will be contained in the body of command.

First, we check an example of ``merkleblock`` command.

Following is the binary data of the body of ``merkleblock`` command:
```ggg
00000020df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bded692d14d4dc7c835b
67d8001ac157e670bfOd00000aba412a0d1480e370173072c9562becffe87aa661c1e4a6dbc305d38ec5dc088a7cf92e6458aca7b32edae818f9c2c98c37e06bf72ae0ce80649a386
55ee1e27d34d9421d940b16732f24b94023e9d572a7f9ab8023434a4feb532d2adfc8c2c2158785d1bd04eb99df2e86c54bc13e139862897217400def5d72c280222c4cbaee7261831
e1550dbb8fa82853e9fe506fc5fda3f7b919d8fe74b6282f92763cef8e625f977af7c8619c32a369b832bc2d051ecd9c73c51e76370ceabd4f25097c256597fa898d404ed53425de608
ac6bfe426f6e2bb457f1c554866eb69dcb8d6bf6f880e9a59b3cd053e6c7060eeacaacf4dac6697dac20e4bd3f38a2ea2543d1ab7953e3430790a9f81e1c67f5b58c825acf46bd02848
384eebe9af917274cdfbb1a28a5d58a23a17977defode10d644258d9c54f886d47d293a411cb6226103b55635
```

Let's put the chunk of binary data above into fields:

1. the first 4 bytes in little endian format is version number: ``0000002``

2. the following 32 bytes in little endian format is id of previous block:
``0df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000``

3. the following 32 bytes in little endian format is value of merkle root:
``ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bded692d14d4``

4. the following four bytes in little endian format is timestamp: ``dc7c835b``

5. the following four bytes is named bits: ``67d8001a``

6. the following four bytes is named nonce: ``c157e670``

7. the following four bytes in little endian format is number of total transactions: ``bfOd0000``

8. the following 1 byte is variant int, the number of hashes: 0a

9. the following chunk of data are hash values of all transactions, its length is 32 * value from step 7:
```ggg
ba412a0d1480e370173072c9562becffe87aa661c1e4a6dbc305d38ec5dc088a7cf92e6458aca7b32edae818f9c2c98c37e06bf72ae0ce80649a386
55ee1e27d34d9421d940b16732f24b94023e9d572a7f9ab8023434a4feb532d2adfc8c2c2158785d1bd04eb99df2e86c54bc13e139862897217400def5d72c280222c4cbaee7261831
e1550dbb8fa82853e9fe506fc5fda3f7b919d8fe74b6282f92763cef8e625f977af7c8619c32a369b832bc2d051ecd9c73c51e76370ceabd4f25097c256597fa898d404ed53425de608
ac6bfe426f6e2bb457f1c554866eb69dcb8d6bf6f880e9a59b3cd053e6c7060eeacaacf4dac6697dac20e4bd3f38a2ea2543d1ab7953e3430790a9f81e1c67f5b58c825acf46bd02848
384eebe9af917274cdfbb1a28a5d58a23a17977defode10d644258d9c54f886d47d293a411cb622610
```

10. the final four bytes are named flag bits: 3b55635

The first six fields are the same as ``getheader`` command, the last four fields are use for proof-of-inclusion.
The value from step 7 is the number of hash values in the list as we mentioned in
the previous sector.


## Bloom Filter
In the previous section,
we ask our full-node peer
to return a ``merkleblock`` command that we can verify whether given transactions of interested(the green boxes)
are included in a block or not.
And most of the time we or the client doesn't want the full-node knows which transactions are interested to us,
therefore, we want to hide our target transactions in a group of transactions (The leafs of the merkle tree).
Then we need an effective method
to transfer info about that group of transactions to full-node.That's where the data structure and algorithm of bloom filter comes into play.

There are 1.041 billion of transactions for bitcoin blockchain now,
how can we quickly select the given several transactions out from 1 billion?
That's where bloom filter comes into play.
Bloom filter is
a kind of data structure used for big data, think about spider of Google crawling web pages,
given a url how the spider knows whether this page is alread saved on the server of google or not.
The way of
doing this is, given the url, we have a group of buckets that are made up of bits, and we have a group of hash functions, the hash functions will hash the given string to the index of a given bucket,
each time we're using a hash function to hash the url to a given bucket, we check the value of the bucket, 
if there is one bucket with the value of 0, then the given page with the url is not saved before,
if all the bucket we visited have the value 1, then we are sure the page is already saved on the server.


There is a possibility of false-positive for bloom filter,
a given page may not save on the server, 
but the bloom filter gives a positive result which means given the url that its page is not saved on the
server, but all the buckets we visited have the value of one.
The possibility of false-positive can be leveraged by enlarging the number of buckets, the more buckets you have,
the less likely you will have
a false positive.


We have bloom filter and create ``filterload`` command to send info about the filter to the full node.
We still need another command name ``getdata`` to request filtered
block from the full node, a filtered block is asking full node to throw transactions to the filter we sent to it and include any transactions that can be matched by the filter(all
buckets have value 1), then put all those filtered transactions into ``merkleblock`` command.

Let's have a look at the payload of ``getdata`` and put it into fields:

020300000030eb2540c41025690160a1014c577061596e32e426b712c7ca00000000000000030000001049847939585b0652fba793661c361223446b6fc41089b8be00000000000000

1. At the beginning its variant, the value in the above data is 0x2 then we only need to get one byte.

2. The following four bytes is type of data item in little endian format: 03000000
(tx: 01000000, block: 02000000, filtered block: 03000000, compact block 04000000)

3. the following 32 bytes are hash identifier:
30eb2540c41025690160a1014c577061596e32e426b712c7ca00000000000000030000001049847939585b0652fba793661c361223446b6fc41089b8be00000000000000

In the payload of ``getdata`` message, if we set the type field to value 3,
then we are asking the full node to return ``merkleblock`` command.


## Segwit
In the previous section of talking about transaction, we have seen some transactions have a bit of segwit set to one.
That indicates such transactions is a kind of "segregated witness" transaction, it is an
upgrade of the traditional transaction, and now it is almost the mainstream transaction.
In this section, we will go into the details of segwit transaction.

There are many benefits brought by segwit transaction compare with the old style transaction:

1. Block size increase

2. Transaction malleability fix

3. segwit versioning for clear upgrade paths

4. quadratic hashing fix

5. offline wallet fee calculation security

The list above is not easy to understand,
we may understand them
after we going to the details of segwit transaction that is pay-to-witness-pubKey-hash transaction
(p2wpkh).
We have seen the
pay-to-pubKey-hash transaction(p2pkh) before, and p2wpkh is an upgrade for p2pkh.
In p2pkh, we combine instructions with data together, but in p2wpkh transaction, we seperate data in ScriptSig to
its own witness field.


There is a jargon word "transaction malleability";
it is the ability to change the transaction id without changing the transaction's intention.
The malleability of transaction ID will brighten many
security breaks for creation of a payment channel,
as we know transaction id is the hash result for content of the transaction,
if any data changed in the transaction and the hash will be invalid.
But it is possible that the scirptSig field changed for the transaction input
may keep the transaction hash remain the same,
because this field will be cleared before computing the transaction
hash.


Therefore, changing the ScriptSig field for transaction input will not affect the transaction hash result.
If the transaction data can be changed without affecting its hash result, then the uniqueness
of a transaction will not be guaranteed by the hash id.
In order to mitigate the problem, bring by emptying ScriptSig field when computing the hash,
p2wpkh transaction will separate data from the scriptSig 
field and put it into another field.


Let's have a look at segwit transaction:
```ggg
010000000115e180dc28a2327e687facc33f10f2a20da717e5548406f7ae8b4c811072f8560100000000ffffffff0100b4f505000000001976a9141d7cd6c75c2e86f4cbf98eaed221b30bd9a0b92888ac00000000
```

1. the first four bytes in little endian format is version: ``01000000``

2. the following field is variant it is the count of input: ``0x01``

3. the following 32 bytes in little endian is previous transaction hash:
``15e180dc28a2327e687facc33f10f2a20da717e5548406f7ae8b4c811072f856``

4. the following four bytes in little endian is previous transaction index: ``01000000``

5. the following one byte 0x00 is scriptSig

6. the following four bytes in little endian is sequence number: ``ffffffff``

7. the following is variant for the count of output: ``0x01``

8. the following 8 bytes in little endian is output amount: ``0b4f50500000000``

9. the following data chunk is ScriptPubKey: ``1976a9141d7cd6c75c2e86f4cbf98eaed221b30bd9a0b92888ac``

10. the last four byte in little endian is lock-time: ``00000000``

Let's check the same transaction with segwit upgrade:

``0100000000010115e180dc28a2327e687facc33f10f2a20da717e5548406f7ae8b4c811072f8560100000000ffffffff0100b4f505000000001976a9141d7cd6c75c2e86f4cbf98eaed221b30bd9a0b92888ac02483045022100df7b7e5cda14ddf91290e02ea10786e03eb11ee36ec02dd862fe9a326bbcb7fd02203f5b4496b667e6e281cc654a2da9e4f08660c620a1051337fa8965f727eb19190121038262a6c6cec93c2d3ecd6c6072efea86d02ff8e3328bbd0242b20af3425990ac00000000``

1. the first four bytes in little endian is version: ``01000000``

2. the following one byte is segwit marker: ``00``

3. the following one byte is segwit flag: ``01``

4. the following field is variant for input count: ``01``

5. the following 32 bytes in little endian is previous transaction hash:
``15e180dc28a2327e687facc33f10f2a20da717e5548406f7ae8b4c811072f856``

6. the following 4 bytes in little endian is previous transaction index: ``01000000``

7. the following one byte is scriptSig: ``00``

8. the following 4 bytes in little endian format is sequence: ``ffffffff``

9. the following variant is number of output: ``01``

10. the following 8 bytes in little endian is output amount: ``00b4f50500000000``

11. the following chunk is scriptPubKey: ``1976a9141d7cd6c75c2e86f4cbf98eaed221b30bd9a0b92888ac``

12. the following data chunk is witness:
``02483045022100df7b7e5cda14ddf91290e02ea10786e03eb11ee36ec02dd862fe9a326bbcb7fd02203f5b4496b667e6e281cc654a2da9e4f08660c620a1051337fa8965f727eb19190121038262a6c6cec93c2d3ecd6c6072efea86d02ff8e3328bbd0242b20af3425990ac``

do the following for each input:
----> 1, the first field is variant, it is a number of items: 0x02
-----> item:
-----> variant, length of item
-----> content of item

13. the last four bytes in little endian is lock-time: 00000000


Compare with p2pkh transaction the p2wpkh has three more fields: segwit marker, segwit flag, and witness.
The field of witness contains tow fields: signature and pubKey.
And the scriptPubKey will contain
two parts, one is instruction OP_0, the second is 20 bytes hash.


when executing the script, the first instruction will push 0 onto the stack,
then a 20-byte hash will push to stack.



For older version of full-node that can not handle segwit transaction,
it will stop here since there is nothing for the script,
and the top element on the stack is not 0 and the result can be seen as
success.
Nodes capable of handing segwit transaction, it will notice the pattern that is OP_0 <20-byte hash>,
it will take the pubKey and signature from witness field and reconstruct the script.


Now we can handle the script as before, and when executing the OP_HASH160, we will put the hash result and the 20 byte hash both on to the stack, and if the OP_EQUALVERIFY will check their match and
the OP_CHECKSIG will check the signature is valid or not, if they are all success, there would be value 1 on the stack.