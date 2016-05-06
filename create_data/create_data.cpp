#include <iostream>
#include <cstdio>
#include <random>

// Dimension of the matrix (height and width are the same)
constexpr uint32_t MATRIX_SIZE=1500;

int main(int argc, char** argv)
{
    FILE* fp = fopen("data.bin", "w");
    if (!fp)
    {
        std::cerr << "Error: could not open data.bin for writing" << std::endl;
        return 1;
    }

    // setup randomizer (0 -> int32_t max value)
    std::random_device rd;
    std::mt19937 gen(rd());
    std::uniform_int_distribution<uint32_t> dis(0, std::numeric_limits<int32_t>::max());

    // calculate size: height * width
    uint32_t size = MATRIX_SIZE * MATRIX_SIZE;

    // write 4 byte unsigned integers
    for (size_t i = 0; i < size; ++i)
    {
        uint32_t v = dis(gen);

        if (fwrite(&v, 4, 1, fp) != 1)
        {
            std::cerr << "Error: failed trying to write random uint32_t to file" << std::endl;
            return 1;
        }
    }

    fclose(fp);
    fp = nullptr;

    return 0;
}

