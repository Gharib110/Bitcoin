# Bitcoin
Bitcoin Implementation in Golang

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

## P2SH Signature Verification
Signature verification is tricky for P2SH transaction.
Since it has several keys and signatures and order of signature should be the same as public key
which means if the order of signature is n then its corresponding
 public key should have order m >= n.

Now comes to how to construct the signature message z. Let's take it step by step,
the following data chunk is a P2SH transaction and the chunk in { and } is the scriptSig of integers input:

```g
0100000001868278ed6ddfb6c1ed3ad5f8181eb0c7a385aa0836f01d5e4789e6bd304d87221a
000000
{
db00483045022100dc92655fe37036f47756db8102e0d7d5e28b3beb83a8fef4f5dc05
59bddfb94e02205a36d4e4e6c7fcd16658c50783e00c341609977aed3ad00937bf4ee942a89937
0148304502210Oda6bee3c93766232079a01639d07fa869598749729ae323eab8eef53577d611b02207bef
15429dcadce2121ea07f233115c6f09034c0be68db99980b9a6c5e75402201475221022626e955ea6ea6d9
8850c994f9107b036b1334f18ca8830bfff1295d21cfdb702103b287eaf122eea69030a0e9feed096bed8045c8b
98bec453e1ffac7fbdbd4bb7152ae
}
ffffffff04d3b11400000000001976a914904a49878c0adfc3aa05de7afad2c
c15f483a56a88ac7f400900000000001976a914418327e3f3dda4cf5b9089325a4b95abdfa0334088ac722c0c0000
0000001976a914ba35042cfe9fc66fd35ac2224eebdafd1028ad2788acdc4ace020000000017a91474d691da1574e
6b3c192ecfb52cc8984ee7b6c568700000000

```

1. find the scriptSig of input replace it with 00 as following:

```g
0100000001868278ed6ddfb6c1ed3ad5f8181eb0c7a385aa0836f01d5e4789e6bd304d87221a
000000
{
00
}
ffffffff04d3b11400000000001976a914904a49878c0adfc3aa05de7afad2c
c15f483a56a88ac7f400900000000001976a914418327e3f3dda4cf5b9089325a4b95abdfa0334088ac722c0c0000
0000001976a914ba35042cfe9fc66fd35ac2224eebdafd1028ad2788acdc4ace020000000017a91474d691da1574e
6b3c192ecfb52cc8984ee7b6c568700000000
```

2. the owner of P2SH transaction has the binary content of the redeem script,
and we use the binary content of the redeem script to replace the 00 above,
for example, the content of the given redeem script is :
```g
475221022626e955ea6ea6d98850c994f9107b036b1334f18ca8830bfff1295d21cfdb702103b287eaf122eea69030a0e9feed096bed8045c8b98bec453e1ffac7fbdbd4bb7152ae
```
Then after replacing the 00 above, we have the following:
```g
0100000001868278ed6ddfb6c1ed3ad5f8181eb0c7a385aa0836f01d5e4789e6bd304d87221a0000000
{
475221022626e955ea6ea6d98850c994f9107b036b1334f18ca8830bfff1295d21cfdb702103b287eaf122eea69030a0e9feed096bed8045c8b98bec453e1ffac7fbdbd4bb7152ae
}
ffffffff04d3b11400000000001976a914904a49878c0adfc3aa05de7afad2c
c15f483a56a88ac7f400900000000001976a914418327e3f3dda4cf5b9089325a4b95abdfa0334088ac722c0c0000
0000001976a914ba35042cfe9fc66fd35ac2224eebdafd1028ad2788acdc4ace020000000017a91474d691da1574e
6b3c192ecfb52cc8984ee7b6c568700000000

```

3. convert the value of hash type SIGHASH_ALL into four bytes and append to the end of the above data:
```g
0100000001868278ed6ddfb6c1ed3ad5f8181eb0c7a385aa0836f01d5e4789e6bd304d87221a000000
{
475221022626e955ea6ea6d98850c994f9107b036b1334f18ca8830bfff1295d21cfdb702103b287eaf122eea69030a0e9feed096bed8045c8b98bec453e1ffac7fbdbd4bb7152ae
}
ffffffff04d3b11400000000001976a914904a49878c0adfc3aa05de7afad2c
c15f483a56a88ac7f400900000000001976a914418327e3f3dda4cf5b9089325a4b95abdfa0334088ac722c0c0000
0000001976a914ba35042cfe9fc66fd35ac2224eebdafd1028ad2788acdc4ace020000000017a91474d691da1574e
6b3c192ecfb52cc8984ee7b6c568700000000[01000000]
```

4. Do hash256 on the above data and the result is the signature message.

## P2SH transaction for multiple signature verification 

For some situation, we may need multiple parties to control the release of fund.
For example, if the fund needs to prove by several board members,
each one has his/her own private key, the fund can only be released
if all board members sign to the contract.

In order to support multiple private keys for one transaction, we need to use the op code with name OP_CHECKMULTISIG,
its hex value is 0xae.
This is a quit complicated opeation; it needs many elements on the stack
to operate, the following binary data is a contract with two public keys:

```g
514104fcf07bb1222f7925f2b7cc15183a40443c578e62ea17100aa
3b44ba66905c95d4980aec4cd2f6eb426d1b1ec45d76724f26901099416b9265b
76ba67c8b0b73d210202be80a0ca69c0e000b97d507f45b98c49f58fec6650b64ff70e6ffccc3e6d0052ae
```
Let's break it down into pieces:

1. the first byte 0x51 is an op code OP_1, which means push value 1 on to the stack

2. the second byte 0x41 is the length for the following data chunk which is the first public key

3. the following 0x41 bytes of data is the raw data of the first publick key:

04fcf07bb1222f7925f2b7cc15183a40443c578e62ea17100aa
3b44ba66905c95d4980aec4cd2f6eb426d1b1ec45d76724f26901099416b9265b
76ba67c8b0b73d

4. the byte following the end of the first public key is 0x21, its length of the second public key

5. the 0x21 bytes following is the second public key:
   0202be80a0ca69c0e000b97d507f45b98c49f58fec6650b64ff70e6ffccc3e6d00

6. The byte following the end of second public key is 0x52, it is an op code OP_2, which means push value 2 on the stack

7. The final byte is an op code OP_CHECKMULTISIG

The above script is scriptPubKey from the output of a previous transaction.
The corresponding scriptSig in current transaction input is as following:
```g
00483045022100e222a0a6816475d85ad28fbeb66e97c
931081076dc9655da3afc6c1d81b43f9802204681f9ea9d52a31c
9c47cf78b71410ecae6188d7c31495f5f1adfeOdf5864a7401
```

1. the first byte 0x00 is op code OP_0, it is used to push an empty array on the stack

2. the second byte 0x48 is the length of signature

3. the data chunk following the second byte is belongs to signature:

```g
3045022100e222a0a6816475d85ad28fbeb66e97c
931081076dc9655da3afc6c1d81b43f9802204681f9ea9d52a31c
9c47cf78b71410ecae6188d7c31495f5f1adfeOdf5864a7401
```

