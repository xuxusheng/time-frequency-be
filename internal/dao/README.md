# dao

全称 data access object

数据接入层，提供数据库相关操作方法，尽量保证每个方法的原子性，让上一层的 service 来进行调度，事务性操作也在上一层中进行编排。

dao 或 repo 之类的划分，有很多概念和方法，但是没太有必要纠结这些。