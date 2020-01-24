## Cache
Certain stim commands use caching to speed up operations.  The structure of the cache is as follows.
```
├── ${STIM_CACHE_PATH}/   # Set via environment variable
│   ├── bin/              # Storage for binary executables
│   │   ├── darwin/       # Versioned MacOS binaries
│   │   ├── linux/        # Versioned Linux binaries
```
