# Bitcoin
Bitcoin Implementation in Golang

## Elliptic Curve
Bitcoin rely heavily on a match object called elliptic curve, without this math structure bitcoin will like a castle on beach, 
it will collapse in any time.What is elliptic curve, it's an equation like this: y^2 = x^3 + ax +b, and its shape just like following:

- For bitcoin, its elliptic curve has a name: secp256k1 and its equation is y^2 = x ^ 3 + 7.

## Signature and Verify
For signature, we can't use SEC for encoding, because there is no close relation between r and s, you know one of them, but you can't derive the other. There is a scheme to encode signature
called DER (Distinguished Encoding Rules), following are steps for encoding signature:

1. set the first byte to 0x30

2. the second byte is the length of signature (usually 0x44 or 0x45)

3. the third byte set to 0x02, this is an indicator that following bytes are for r.

4. transfer r into bytes array, if the first byte of r is >= 0x80, than we append a byte 0x00 in the beginning of the bytes, compute the length of the bytes array, append the value of length
after the marker byte 0x02 in step 3, and append the bytes array following the value of length

5. add a marker byte 0x02 at the end of bytes array from step 4.

6. append s as how we append r in step 4.

We need to encode the length of r and s, because r and s at most has 32 bytes, but some times their length may shorter than this. Let's see an example for the encoding of r and s:

30 45 02 21 00 ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f 02 20 7a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed

the first byte is 0x30 as mentioned above, the second byte is 0x45 that is the total length of r and s. The third byte is marker 0x02 indicating from next byte is the beginning of bytes array for r.
According to 3, the byte following from marker 0x02 is the length of bytes array for r, the value is 0x21. Byte following 0x21 is 0x00, by step 4, it indicates the first byte of r is >= 0x80, we can
see the byte follow 0x00 is 0xed which is indeed more than 0x80, the length of r is 0x21 - 1 = 0x20, which means the following 32 bytes are bytes array for r, we extract it out as following:

r: 81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f

following the last byte of r is 0x20 , it is indicator for the beginning of s, the byte following 0x02 is 0x20 which indicates the length of s is 0x20, the following byte is not 0x00, which means the
first byte of s is not more than 0x80, and the byte following 0x20 is 0x7a , this is the beginning byte for s and is smaller than 0x80, the byte array for s is :

s: 7a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed

## Wallet Address
As a bitcoin user, we always need to send or receive bitcoins from other, this will require to let others know your wallet address. 
Because wallet address need to be read by human, all the encoding schema
we have before are produce result in binary, therefore we need another scheme to create wallet address in human friendly way.

Wallet address is actually generated from public key, and it needs to satisfy the following requirement:

1. easy to read and write, user may want to memorize it or write it down on paper

2. Not too long for sending over the internet

3. It should be secure, and harder to make mistake, you don't want you fund transfer to people unknown to you!

The base58 encoding scheme can help us to achieve three goals. Compare with the commonly use of base64, it removes characters like l , I, 0, O, -, _ because they are easy to confuse with each other.
Because the encoding schema uses all numbers, and uppercase and lowercase letters and remove 0 O, l , I, which means it will use 58 characters in the encoding process.


## Transaction Validation

For a bitcoin node, one of its major task is to validate a transaction, there are several steps to take for it, the first thing is to check the output can match to the transaction. For example if a transaction
is about "jim using 10 dollars to by a cup of coffee with price of 3 dollars", then we need to check :

1, jim really has 10 dollars

2, the amount left after buying the coffee should be 7 dollars

If the transaction is honest, then the input of the transaction(10 dollars) should greater than the output of the transaction(7 dollars), that is when we use the amount of input minus the amount of the output
the result should be positive, if the result is negative, then the transaction is "dishonest" it want to fake money from air.