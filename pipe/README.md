## shunt

- 用于多路串行化处理，基本原理是启动一系列go routine，在任务分发时通过计算slot，
将不同的请求分发到不同slot对应的go routine上执行。
