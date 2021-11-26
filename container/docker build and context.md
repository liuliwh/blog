- [What is the docker image?](#what-is-the-docker-image)
- [What is the build context?](#what-is-the-build-context)
- [Why need build context?](#why-need-build-context)
  - [Two hosts are involved:](#two-hosts-are-involved)
- [How to speicfy the build context?](#how-to-speicfy-the-build-context)
  - [Git repo as build context](#git-repo-as-build-context)
    - [Summary](#summary)
  - [Remote tarball context](#remote-tarball-context)
    - [Summary](#summary-1)
  - [No build context](#no-build-context)
    - [This implies **a Dockerfile that uses COPY or ADD will fail if this syntax is used**](#this-implies-a-dockerfile-that-uses-copy-or-add-will-fail-if-this-syntax-is-used)
    - [Example](#example)
  - [Use .dockerignore exlude files sending to daemon](#use-dockerignore-exlude-files-sending-to-daemon)
- [Dockerfile](#dockerfile)
  - [ARG vs ENV](#arg-vs-env)
- [Relationship between build context and Dockerfile](#relationship-between-build-context-and-dockerfile)
- [Others](#others)
  - [Set build-time variables (--build-arg)](#set-build-time-variables---build-arg)
- [References](#references)

# What is the docker image?
The image is loose collection of *independent layers*, it is just *a configuration file that lists the layers and some metadata*. 
- The layers are where the data lives, each layer has no concept of being part of an overall bigger image. Each layer has distribution ID(compressed then digest).
  
```bash
$ docker image inspect --format "{{ .RootFS.Layers }}" 164a0f0bc8c2 
[sha256:6acbb9f1f55dae18235a3b234c0829d0ec790f8eaa1599834b17d2c9348a1f16 sha256:43656e43813a6559e595b5f0bbf758229b5fd6de80a0c6f439db880294469d41 sha256:825abb30c5e28c3f5de94a4af2982dfb350a52ef642145a5a29bab477bd19125]
$ docker image history 164a0f0bc8c2
IMAGE          CREATED        CREATED BY                                      SIZE      COMMENT
164a0f0bc8c2   2 weeks ago    /bin/sh -c #(nop)  ENTRYPOINT ["/httpserver"]   0B        
<missing>      2 weeks ago    /bin/sh -c #(nop)  USER nonroot:nonroot         0B        
<missing>      2 weeks ago    /bin/sh -c #(nop)  EXPOSE 8080                  0B        
<missing>      2 weeks ago    /bin/sh -c #(nop) COPY file:9e9ec92181241bf4…   7.21MB    
<missing>      2 weeks ago    /bin/sh -c #(nop) WORKDIR /                     0B        
<missing>      2 weeks ago    /bin/sh -c #(nop)  ENV VERSION=v1.0             0B        
<missing>      51 years ago   bazel build ...                                 17.4MB    
<missing>      51 years ago   bazel build ...                                 1.82MB   
```
# What is the build context?
```bash
docker build [OPTIONS] PATH|URL|-
Sending build context to Docker daemon  2.607kB
```
A build’s context is the set of files located in the specified PATH or URL of *docker build* command.The directories and files which will be avaible to the docker engine.
# Why need build context?
Docker CLI and Docker engine might not be running on the same machine, so when issue 'docker build', the CLI sends the build context to the Docker engine, the build process can refer to any files in the context. 
## Two hosts are involved:
1. Docker CLI host. The host where you run the docker build command.
2. Docker daemon host. The host the Docker daemon is running on. 
> It is not necessarily the same host from which the build command is being issued 
# How to speicfy the build context?
## Git repo as build context
```bash
# simple example, build context used refs/heads/master/
docker build https://github.com/docker/rootfs.git
# complicate example, build context used refs/tags/tag/folder
docker build https://github.com/docker/rootfs.git#tagorbranch:folder
```
> Under the hood <br>Docker performs **git clone --recursive on the local machine**, and **sends those files as build context to the daemon**. It requires git to be installed on the host where you run the docker build command.

### Summary 
**Docker CLI host downloads and sends** the build context to Docker daemon host
## Remote tarball context
### Summary 
**Docker deamon host downloads and untar** the build context.
## No build context
> If you use STDIN or specify a URL pointing to a plain text file, the system places the contents into a file called Dockerfile, and any -f, --file option is ignored. In this scenario, there is no context.
1. Remote plaintext
2. STDIN (-)
### This implies **a Dockerfile that uses COPY or ADD will fail if this syntax is used**
### Example
To pipe a Dockerfile from STDIN
```bash
 docker build - < Dockerfile
```
pipe Dockerfile through stdin for one-off build. For example. 
```bash
vagrant@cncamp:~/test$ # echo -e means backslash escape 
vagrant@cncamp:~/test$ echo -e 'FROM busybox\nRUN echo "hello world"' | docker build -
Sending build context to Docker daemon  2.048kB
Step 1/2 : FROM busybox
 ---> 16ea53ea7c65
Step 2/2 : RUN echo "hello world"
 ---> Using cache
 ---> d7aa7a68e127
Successfully built d7aa7a68e127
```
## Use .dockerignore exlude files sending to daemon
Since the docker context will be sent to docker daemon, build context -> image size -> build time -> pull/push container size. so we should avoid sending big build context. 
From docker doc
> You can even use the .dockerignore file to exclude the Dockerfile and .dockerignore files. These files are still sent to the daemon because it needs them to do its job. But the ADD and COPY instructions do not copy them to the image.
```bash
Dockerfile*
*.md
!README.md
```
# Dockerfile
Specify a Dockerfile with option (-f)
## ARG vs ENV
|     | Persist in image and container| Impact on build caching | Priority | Scope                                                                                                                       | Note                                                                           |
|-----|---------------|-------|----------|-----------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------|
| ENV | Y             | N     | H        |                                                                                                                             |                                                                                |
| ARG | N             | Y     | L        | from the line which is defined. <br>end of the build stage where it was defined. | To use an arg in multiple stages, each stage must include the ARG instruction. |
|     |               |       |          |                                                                                                                             |                                                                                |
# Relationship between build context and Dockerfile
1. By default the docker build command will look for a **Dockerfile at the root of the build context.**
2. The build process can refer to any of the files in the context. so, COPY/ADD/... instructions can reference the files in the context.
# Others
## Set build-time variables (--build-arg)

# References
1. https://docs.docker.com/engine/reference/commandline/build/
2. https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#understand-build-context
3. https://docs.docker.com/engine/reference/builder/#dockerignore-file