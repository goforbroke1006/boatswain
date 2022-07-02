cmake_minimum_required(VERSION 3.16)
project(boatswain)

set(CMAKE_CXX_STANDARD 14)

add_executable(boatswain main.cpp)

#find_package(CryptoPP REQUIRED)
target_link_libraries(boatswain pthread cryptopp)


enable_testing()
add_subdirectory(tests)
