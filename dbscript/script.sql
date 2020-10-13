-- Title: 聊天室資料格式腳本
-- Author: Terry Huang

DROP DATABASE IF EXISTS chatroom;
CREATE DATABASE chatroom;
USE chatroom;

-- 創建使用者資料表
DROP TABLE IF EXISTS Users;
CREATE TABLE Users (
    userId varchar(30) primary key, -- 使用者Id
    userName varchar(64),           -- 使用者名稱
    password char(64) not null,     -- 使用者密碼
    createTime datetime not null    -- 創建的時間
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 創建聊天室資料表
DROP TABLE IF EXISTS Hubs;
CREATE TABLE Hubs (
    hubId varchar(30) primary key, -- 聊天室Id
    hubName varchar(64),           -- 聊天室名稱
    userId varchar(30),            -- 創建的使用者
    createTime datetime not null,  -- 創建的時間
    FOREIGN KEY (userId) REFERENCES Users( userId )
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 創建聊天室與使用者關聯表
DROP TABLE IF EXISTS Registers;
CREATE TABLE Registers (
    hubId varchar(30) not null,     -- 聊天室Id
    userId varchar(30) not null,    -- 使用者Id
    type tinyint not null,          -- 身份: 0 管理員 1 一般使用者
    registerTime datetime not null, -- 加入聊天室的時間
    FOREIGN KEY ( hubId ) REFERENCES Hubs( hubId ),
    FOREIGN KEY ( userId ) REFERENCES Users( userId )
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 創建訊息資料表
DROP TABLE IF EXISTS Messages;
CREATE TABLE Messages (
    hubId varchar(30) not null,     -- 聊天室Id
    messageId int,                  -- 訊息編號
    userId varchar(30) not null,    -- 使用者Id
    content varchar(255) not null,  -- 訊息內容
    createTime datetime not null,   -- 發送訊息的時間
    FOREIGN KEY ( hubId ) REFERENCES Hubs( hubId ),
    FOREIGN KEY ( userId ) REFERENCES Users( userId )
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
