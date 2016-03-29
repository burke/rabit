#include <stdint.h>
#include <stdlib.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <errno.h>

const uint32_t windowSize = 64;
const uint32_t charOffset = 31;

const uint32_t blobBits = 17;
const uint32_t blobSize = 1 << blobBits; // 128k
const uint32_t blobMask = blobSize - 1;
const uint32_t splitMask = -1U & blobMask;

const uint32_t maxBlobSize = 1 << 20;
const uint32_t tooSmallThreshold = 64 << 10;

struct RollSum {
  uint32_t s1;
  uint32_t s2;
  uint8_t window[windowSize];
  int wofs;
};

struct RollSum *
NewRollSum()
{
  struct RollSum * rs = calloc(1, sizeof (struct RollSum));
  rs->s1 = windowSize * charOffset;
  rs->s2 = windowSize * (windowSize - 1) * charOffset;
  return rs;
}

#define Roll(rs, add) do { \
  uint32_t drop = (rs)->window[(rs)->wofs]; \
  (rs)->s1 += (add) - drop; \
  (rs)->s2 += (rs)->s1 - windowSize*(drop+charOffset); \
  (rs)->window[(rs)->wofs] = (add); \
  (rs)->wofs = ((rs)->wofs + 1) % windowSize; \
} while (0)

int main()
{
  struct RollSum * rs;
  int fd, i, blobSize;
  char *data;
  struct stat sbuf;
  bool onSplit;

  if ((fd = open("/opt/railgun/share/images/railgun-common-services-0.1.4.img", O_RDONLY)) == -1) {
    perror("open");
    exit(1);
  }

  if (stat("/opt/railgun/share/images/railgun-common-services-0.1.4.img", &sbuf) == -1) {
    perror("stat");
    exit(1);
  }

  data = mmap((caddr_t)0, sbuf.st_size, PROT_READ, MAP_SHARED, fd, 0);
  if (data == (caddr_t)(-1)) {
      perror("mmap");
      exit(1);
  }

  rs = NewRollSum();

  blobSize = 0;
  for (i = 0; i < sbuf.st_size; i++) {
    blobSize++;
    Roll(rs, data[i]);
    onSplit = ((rs)->s2 & blobMask) == splitMask;
    if (blobSize == maxBlobSize || (onSplit && blobSize > tooSmallThreshold)) {
      printf("%d\n", i);
      blobSize = 0;
    }
  }
  free(rs);
  return 0;
}
