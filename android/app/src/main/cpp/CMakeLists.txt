cmake_minimum_required(VERSION 3.22.1)
project(honoka_android)

add_library(
    honoka_android
    SHARED
    honoka_android.cpp
)

find_library(log-lib log)

find_library(dl-lib dl)

target_link_libraries(
    honoka_android
    ${log-lib}
    ${dl-lib}
)
