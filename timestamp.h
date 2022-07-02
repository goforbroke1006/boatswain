//
// Created by goforbroke on 02.07.22.
//

#ifndef BOATSWAIN_TIMESTAMP_H
#define BOATSWAIN_TIMESTAMP_H

#include <chrono>

uint64_t getTimestamp() {
    const auto now = std::chrono::system_clock::now();
    return std::chrono::duration_cast<std::chrono::seconds>(now.time_since_epoch()).count();
}

#endif //BOATSWAIN_TIMESTAMP_H
