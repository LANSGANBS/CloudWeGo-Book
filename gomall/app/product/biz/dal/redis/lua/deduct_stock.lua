-- deduct_stock.lua
-- 原子性扣减库存 Lua 脚本
-- KEYS[1]: 库存键名 (stock:product:{product_id})
-- ARGV[1]: 扣减数量
-- ARGV[2]: 当前时间戳（用于日志）
-- 返回值: 
--   正数: 扣减后的剩余库存
--   -1: 库存不足
--   -2: 库存键不存在

local stock_key = KEYS[1]
local deduct_amount = tonumber(ARGV[1])
local current_stock = tonumber(redis.call('GET', stock_key))

if current_stock == nil then
    return -2
end

if current_stock < deduct_amount then
    return -1
end

local new_stock = redis.call('DECRBY', stock_key, deduct_amount)
return new_stock
