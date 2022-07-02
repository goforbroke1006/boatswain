//
// Created by goforbroke on 02.07.22.
//

#ifndef BOATSWAIN_BLOCKCHAIN_H
#define BOATSWAIN_BLOCKCHAIN_H

#include <string>
#include <utility>
#include <vector>
#include <ostream>

#include "hash.h"
#include "timestamp.h"

class Block {
public:
    explicit Block(uint64_t index, const std::string &previousHash, uint64_t timestamp, const std::string &data)
            : mIndex(index),
              mPreviousHash(previousHash),
              mTimestamp(timestamp),
              mData(data) {
        this->mHash = SHA256HashString("" + std::to_string(index) + previousHash + std::to_string(timestamp) + data);
    }

    uint64_t getIndex() const {
        return mIndex;
    }

    const std::string &getHash() const {
        return mHash;
    }

    void setHash(const std::string &hash) {
        Block::mHash = hash;
    }

    const std::string &getPreviousHash() const {
        return mPreviousHash;
    }

    void setPreviousHash(const std::string &previousHash) {
        Block::mPreviousHash = previousHash;
    }

    uint64_t getTimestamp() const {
        return mTimestamp;
    }

    const std::string &getData() const {
        return mData;
    }

private:
    uint64_t mIndex;
    std::string mHash;
    std::string mPreviousHash;
    uint64_t mTimestamp;
    std::string mData;
};

const char *lockUTF8Octal = "\360\237\224\222";               // https://www.unicodepedia.com/unicode/miscellaneous-symbols-and-pictographs/1f512/lock/
const char *clockFaceOneOClockUTF8Octal = "\360\237\225\220"; // https://www.unicodepedia.com/unicode/miscellaneous-symbols-and-pictographs/1f550/clock-face-one-oclock/
const char *clipboardUTF8Octal = "\360\237\223\213";          // https://www.unicodepedia.com/unicode/miscellaneous-symbols-and-pictographs/1f4cb/clipboard/

std::ostream &operator<<(std::ostream &out, Block *block) {
    out << "# " << block->getIndex() << " "
        << lockUTF8Octal << ": " << block->getHash().c_str() << " "
        << clockFaceOneOClockUTF8Octal << ": " << block->getTimestamp() << " "
        << clipboardUTF8Octal << ": " << block->getData().c_str();
    return out;
}

Block *getGenesisBlock() {
    const uint64_t t_2022_07_02_00_00_0300 = 1656709200;
    return new Block(0, "0", t_2022_07_02_00_00_0300, "Initial Block in the Chain");
}

class BlockChain {
public:
    void startGenesisBlock() {
        mChain.push_back(getGenesisBlock());
    }

    Block *obtainLatestBlock() {
        return mChain.back();
    }

    void generateNewBlock(uint64_t timestamp, const std::string &data) {
        const std::string &prevHash = this->obtainLatestBlock()->getHash();
        auto *pBlock = new Block(this->mChain.size(), prevHash, timestamp, data);
        mChain.push_back(pBlock);
    }

    const std::vector<Block *> &getChain() const {
        return mChain;
    }

    virtual ~BlockChain() {
        for (auto *pBlock: mChain) {
            delete pBlock;
        }
    }

private:
    std::vector<Block *> mChain;
};


#endif //BOATSWAIN_BLOCKCHAIN_H
