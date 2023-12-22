---
title: "{{ replace .Name "-" " " | title }}"
date: {{ .Date }}
draft: false
description: 
resources:
  - src: foo.jpg
    title: Foo
    params:
      authors:
      source:
---
