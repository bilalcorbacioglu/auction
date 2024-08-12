# FairBlock: Revolutionising Blockchain Encryption and Auction Systems

Although blockchain technology offers the potential for decentralisation, security and transparency, these benefits can only be fully realised if certain fundamental issues are addressed. It is essential that critical processes such as auctions are conducted fairly, securely and transparently in order to fully harness the potential of blockchain. This is where FairBlock comes in, offering a revolutionary solution: programmable encryption.

# A Closer Look at Centralised Auctions

A centralised auction system is one in which a single authority or platform exercises control over the process. Participants submit their bids to this central entity, which then determines the winning bid. In such systems, bids are collected, processed, and results are announced from a centralised authority. Even within the context of blockchain technology, many auction processes are still conducted in a centralised manner, which gives rise to significant issues in terms of transparency, security, and fairness.

### Examples and Issues

The opacity of Google's ad auctions, which allocate digital advertising space, renders them susceptible to manipulation, as evidenced by legal actions such as the U.S. Department of Justice's lawsuit against Google. In financial markets, auctions for bond and stock issuances carry the risk of insider trading and manipulation, which can undermine market integrity. Carbon credit auctions give rise to concerns about market manipulation by major players, which could compromise environmental sustainability goals. Furthermore, even blockchain-based auctions like NFT sales, when conducted through centralised platforms, are vulnerable to manipulation and lack transparency. These examples demonstrate the fundamental issues of centralised auctions, which are rooted in the lack of transparency, security, and fairness.

# Looking for a Solution? "Haven’t You Heard of FairBlock?"

FairBlock's technology is set to revolutionize second-price sealed-bid auctions. Participants encrypt their bids and submit them to the blockchain, and once the auction period ends, validators on the network decrypt the bids, ensuring a fair and secure process. Let's break down how this works step by step:


### 1\. **Fairyring**: The Platform for Secure and Fair Transactions

The first step in ensuring fairness and security in auctions is for participants to encrypt their bids. This is where Fairyring comes into play. Fairyring is a specialized chain developed to securely encrypt, decrypt, and execute blockchain transactions. Participants use **Threshold Identity-Based Encryption** (tIBE) to encrypt their bids. This encryption method ensures that transactions can only be decrypted under specific conditions, such as reaching a certain block height.

![Image1](https://raw.githubusercontent.com/bilalcorbacioglu/auction/main/blog/images/d1.png)

### 2\. Secure Storage of Encrypted Bids: Say Hello to **x/pep**

Encrypted bids are sent to the x/pep repository, where they are securely stored until they are ready to be processed at a specific block height. This repository holds the encrypted transactions securely until the predefined conditions are met, at which point they are automatically decrypted and executed by the Fairyring chain.

![Image1](https://raw.githubusercontent.com/bilalcorbacioglu/auction/main/blog/images/t1.png)


### 3\. Keyshare Aggregation: Secure Decryption with Threshold Cryptography

When the auction period ends, the bids must be decrypted. This is where **Threshold Cryptography** comes into play. FairBlock generates a **Master Secret Key (MSK)** during each epoch, which is then split into keyshares distributed among validators. Upon reaching the designated block height, validators submit their shares to derive a unique private key for that block.

To ensure security, the **2/3 \+ 1 rule** is applied. According to this rule, at least two-thirds plus one of the validators must submit their keyshares to derive the block’s private key. 

### 4\. Decrypting and Finalising Bids

When the Fairyring chain reaches the target block height, it derives the necessary keys, as outlined in Step 3, and decrypts the bids. These decrypted bids are processed before any other transactions, ensuring a fair and secure auction outcome.

### Overview

Here's a diagram visualising how these steps work, showing the end-to-end process of encrypting, storing, decrypting, and executing transactions:

![Image1](https://raw.githubusercontent.com/bilalcorbacioglu/auction/main/blog/images/d2.png)

# But Wait, There's More\!

FairBlock's security and transparency solutions are further enhanced by advanced modules and products. In addition to the products mentioned above, FairBlock increases the security and efficiency of critical processes in blockchain networks, ranging from auctions to data security, with the following tools:

**Keyshare Module:** FairBlock’s Keyshare Module enables secure distribution and aggregation of keyshares among validators. This module utilizes threshold cryptography to securely generate and decrypt private keys. Upon reaching a specific block height, validators submit their keyshares, which are then aggregated to derive the necessary private key for that block.

**Share Generator:** FairBlock’s Share Generator product generates the **Master Secret Key (MSK)** during each epoch and divides it into keyshares for distribution among validators. The Share Generator securely creates these keyshares and facilitates their submission to the blockchain.

![Image1](https://raw.githubusercontent.com/bilalcorbacioglu/auction/main/blog/images/t2.png)

**Fairyport:** Fairyport is FairBlock’s **cross-chain** communication module. It securely transfers the private key necessary for decryption to the target chain once the block height conditions are met. Fairyport establishes a secure connection between the Fairyring chain and the target chain, ensuring that transactions are executed fairly and securely.

# A Secure and Fair Future is Possible

FairBlock not only addresses current issues but also lays the groundwork for future needs. For example, the integration of advanced techniques like **Fully Homomorphic Encryption (FHE)** allows data to be processed while still encrypted. This is critical for applications requiring high privacy, such as secure voting systems, confidential data analysis, and beyond.

By offering programmable encryption solutions for decentralized auctions and blockchain applications, FairBlock is building the financial and technological infrastructure of the future. Supported by advanced cryptographic techniques, these solutions maximise security, privacy, and fairness. With robust APIs and flexible integration options for developers, this technology is poised to transform not only the blockchain but also various critical processes in everyday life.