//
// Created by goforbroke on 02.07.22.
//

#ifndef BOATSWAIN_HASH_H
#define BOATSWAIN_HASH_H

#include <string>
#include <crypto++/sha.h>
#include <crypto++/filters.h>
#include <crypto++/hex.h>

/**
 * SHA256HashString get sha256 mHash for provided string.
 * Based on https://cpp.hotexamples.com/ru/site/redirect?url=https%3A%2F%2Fgithub.com%2FAhmedWaly%2FCastOnly/blob/master/OTSL-master/tools.cpp sample.
 * @param plainText context
 * @return sha256 result as a string
 */
std::string SHA256HashString(const std::string &plainText) {
    CryptoPP::SHA256 hash;
    byte digest[CryptoPP::SHA256::DIGESTSIZE];
    hash.CalculateDigest(digest, (const byte*)plainText.c_str(), plainText.length());

    // Crypto++ HexEncoder object
    CryptoPP::HexEncoder encoder;

    // Our output
    std::string output;

    // Drop internal hex encoder and use this, returns uppercase by default
    encoder.Attach(new CryptoPP::StringSink(output));
    encoder.Put(digest, sizeof(digest));
    encoder.MessageEnd();

    // Convert to lowercase if needed
    std::transform(output.begin(), output.end(), output.begin(),
                   [](unsigned char c){ return std::tolower(c); });

    return output;
}

#endif //BOATSWAIN_HASH_H
