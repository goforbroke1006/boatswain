//
// Created by goforbroke on 02.07.22.
//

#include <gtest/gtest.h>

#include "hash.h"

TEST(SHA256HashString_fn, Positive_GenesisHash) {
    std::string phrase = "001656709200Initial Block in the Chain";
    std::string expected = "aed11b71d6c952f7d45756eae9c951f50d2f8902f736662421eb139855f87edd";
    std::string actual = SHA256HashString(phrase);
    ASSERT_EQ(expected, actual);
}