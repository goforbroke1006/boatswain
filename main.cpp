#include <iostream>

#include "blockchain.h"

int main() {
//    std::locale::global(std::locale("en_US.UTF-8"));
//    setlocale(LC_ALL, "");
//    setlocale(LC_CTYPE,"");

    const uint64_t t_2022_07_16_54_04_gmt = 1656780844;
    const uint64_t t_2022_07_16_56_03_gmt = 1656780963;
    const uint64_t t_2022_07_17_02_55_gmt = 1656781375;

    BlockChain blockChain;
    blockChain.startGenesisBlock();
    blockChain.generateNewBlock(t_2022_07_16_54_04_gmt, "Sat Jul 02 2022 16:54:04 GMT+0000");
    blockChain.generateNewBlock(t_2022_07_16_56_03_gmt, "Sat Jul 02 2022 16:56:03 GMT+0000");
    blockChain.generateNewBlock(t_2022_07_17_02_55_gmt, "Sat Jul 02 2022 17:02:55 GMT+0000");

    for (auto *pBlock: blockChain.getChain()) {
        std::cout << pBlock << std::endl;
    }
    return 0;
}
