project(${CMAKE_PROJECT_NAME}-tests)

# Locate GTest
find_package(GTest REQUIRED)
include_directories(${GTEST_INCLUDE_DIRS} ${gtest_SOURCE_DIR})

FILE(GLOB GLOB_TEST_CPP ./../*.cpp ./../*/*.cpp ./../*/*/*.cpp)
message(${GLOB_TEST_CPP})
list(FILTER GLOB_TEST_CPP EXCLUDE REGEX "main.cpp")

add_executable(${PROJECT_NAME} ./main.cpp ${GLOB_TEST_CPP})
add_test(AllTests ${PROJECT_NAME})

target_link_libraries(${PROJECT_NAME} ${GTEST_LIBRARIES} pthread cryptopp)
