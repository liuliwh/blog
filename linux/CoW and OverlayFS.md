- [Linux filesystem](#linux-filesystem)
  - [Bootfs](#bootfs)
  - [Rootfs](#rootfs)
- [Copy-on-Write strategy](#copy-on-write-strategy)
  - [Usage](#usage)
- [Lab - Create an Overlay Filesystem (fuse-overlayfs)](#lab---create-an-overlay-filesystem-fuse-overlayfs)
- [Lab - docker in action](#lab---docker-in-action)
# Linux filesystem
## Bootfs
Bootloader load kernal into memory, then kernal umount bootfs
## Rootfs
Linux set rootfs to Readonly -> check -> readwrite

**Docker set rootfs to Readonly -> check -> union mount a readwrite FS to readonly roootfs, stack CoW layer** 
# Copy-on-Write strategy
Any files of lower layer does not change do not get copied to upper writable layer. This means that the writable layer is as small as possible.
- search from newest(upper) layer down to base layer one layer at a time. When results are found, they are added to a cache to speed future operation.
- perform *copy_up* operation on the first copy of the file that is found, to copy the file to the upper upper writable layer.
- any modifications are made to this copy of the file, and the upper layer cannot see the read-only copy of the file that exists in the lower layer.
## Usage
It is suitable for thin writable layer. Not suitable for heavy-write application.
# Lab - Create an Overlay Filesystem (fuse-overlayfs)
```bash
sh1 $mkdir overlaytest
sh1 $cd overlaytest/
sh1 $mkdir UPPER LOWER WORK MERGED
sh1 $echo "from lower" > ./LOWER/lower.txt
sh1 $echo "from upper" > ./UPPER/upper.txt
sh1 $echo "from lower" > ./LOWER/both.txt
sh1 $echo "from upper" > ./UPPER/both.txt
sh1 $fuse-overlayfs -olowerdir=LOWER,upperdir=UPPER,workdir=WORK MERGED
sh1 $tree
.
├── LOWER
│   ├── both.txt
│   └── lower.txt
├── MERGED
│   ├── both.txt
│   ├── lower.txt
│   └── upper.txt
├── UPPER
│   ├── both.txt
│   └── upper.txt
└── WORK
    └── work

5 directories, 7 files
sh1 $cat MERGED/both.txt 
from upper
```
# Lab - docker in action
```bash
docker inspect --format "{{json .GraphDriver}}" 16aaafe6cb87
{"Data":{"LowerDir":"/var/lib/docker/overlay2/fd35df13f290c726c7609222b8644253ed8ef6b3eff39c87946c12a614ddc42b-init/diff:/var/lib/docker/overlay2/d69cbf65ce321ac59fc836e2cc2015004aed29d6df81256896491deb475ee3a1/diff:/var/lib/docker/overlay2/46b502bafc43251035b4b5d00a5ff525d76cf02899371373d9fe0b803b3511d7/diff:/var/lib/docker/overlay2/704468256bdac88015848bdf23392576c30a269342873e9ae7da6cb96f6a2f66/diff","MergedDir":"/var/lib/docker/overlay2/fd35df13f290c726c7609222b8644253ed8ef6b3eff39c87946c12a614ddc42b/merged","UpperDir":"/var/lib/docker/overlay2/fd35df13f290c726c7609222b8644253ed8ef6b3eff39c87946c12a614ddc42b/diff","WorkDir":"/var/lib/docker/overlay2/fd35df13f290c726c7609222b8644253ed8ef6b3eff39c87946c12a614ddc42b/work"},"Name":"overlay2"}
```