title: my vim golang programming environment
date: 2021-02-14 23:35:19
categories: 技术
tags: [golang,vim,docker]
---

If you are also suffered from vim emulation or facked vim such as ideavim or vscodevim. Today we will show you a new golang development environment under docker or use it out of the box. And easy to tweak by our public dockerfile and built image from docker hub public image. Here's a list of features in the following article.

## features

Autocomplete by language server without import package we need. While showing function parameters and documents.

![image](https://user-images.githubusercontent.com/3077762/107909035-53657b00-6f92-11eb-8b09-49741eb77d18.png)

We use onedark color theme as the basic vim color solution. Here we give two everyday operations examples as below.

1. Rename variable
2. Add tags to struct

[![asciicast](https://asciinema.org/a/qX4GNuz8InY1aqPLRwGMrVRt0.svg)](https://asciinema.org/a/qX4GNuz8InY1aqPLRwGMrVRt0)

The example below use [ETCD](https://github.com/etcd-io/etcd) as demo project show the following feature:

1. tini for container init process
2. auto complete
3. go to definition
4. FZF search file

[![asciicast](https://asciinema.org/a/egDeADL1ITctljs1C9pFjqghG.svg)](https://asciinema.org/a/egDeADL1ITctljs1C9pFjqghG)

Of course gopls (which is a golang toolchain we depend on) use over 40 seconds to scan so large project. But gopls also have a cache feature accelerate open project in the next time.

## how to use it

pull the docker image by command:

```
docker pull lijianying10/golangdev:21Feb7-01
```

docker run by following command

```
docker run -it --rm -v $PWD/etcd:/root/etcd lijianying10/golangdev:21Feb7-01 /bin/bash
```

Attention: alter the dir mapping to your project path, and we highly recommend using `gomod` as the project manager.

## Dockerfile 

ref link: https://github.com/lijianying10/FixLinux/blob/master/golangdev/Dockerfile

## vim dot file

ref link: https://github.com/lijianying10/FixLinux/blob/master/dotfile/.vimrc

### shortcut keys (key maps)

start from [LOC 79](https://github.com/lijianying10/FixLinux/blob/master/dotfile/.vimrc#L97) of dot file.

``` viml
nmap <M-p> :TagbarToggle<CR> " view tag bar
imap <M-p> <esc>:TagbarToggle<CR>i
nmap <M-u> :NERDTreeToggle<CR> " view file list
imap <M-u> <esc>:NERDTreeToggle<CR>
nmap <C-c> :q<CR> " exit 
nmap <M-o> :tabn<CR> " tab next
imap <M-o> <esc>:tabn<CR>
nmap <M-i> :tabp<CR> " tab previous
imap <M-i> <esc>:tabp<CR>
nmap <M-l> :w<CR>:GoMetaLinter<CR> " linter 
nmap <M-n> <Plug>(coc-definition) " go to definition
nmap <C-z> :undo<CR> " undo
nmap <M-y> :GoErrCheck<CR> " go error check
nmap <C-s> :w<CR> " save
imap <C-s> <esc>:w<CR>
imap <M-c> <esc>:pc<CR>
nmap <M-c> :pc<CR> " close preview window
nmap <leader>r :Ack<space> " search hole project document: https://github.com/mileszs/ack.vim
nmap <leader>t :FZF<CR> " zfz file search
```

example key mapping:

1. M-p means `Meta + p` Option key for mac and alt key for windows keyboard
1. C-s means `Ctrl + s` 
1. `<leader>t` means press `\` and then press t

