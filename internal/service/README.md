# service

service 服务层，对上供 controller 层调用，对下对 dao 层进行调度。

尽量将所有的业务逻辑，都收敛在 service 层中，这样可以使 controller 和 dao 层尽量与业务无关，方便复用及多人协作时分工 

service 不对输入进行格式性的校验（字符串长度、数字大小等等），仅进行业务性校验（名称重复等等）。