-- KEYS[1]: Token bucket key (e.g., "{clientA:user1}:tokens")
-- KEYS[2]: Timestamp key (e.g., "{clientA:user1}:ts")

-- ARGV[1]: Limit (Maximum token capacity)
-- ARGV[2]: Windowms (Time to completely refill the bucket in ms)
-- ARGV[3]: Current timestamp in ms (Passed from Go)

local tokens_key = KEYS[1]
local ts_key = KEYS[2]

local capacity = tonumber(ARGV[1])
local window_ms = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

-- 1. Get current data
local last_tokens = tonumber(redis.call("GET", tokens_key))
if last_tokens == nil then last_tokens = capacity end

local last_ts = tonumber(redis.call("GET", ts_key))
if last_ts == nil then last_ts = now end

-- 2. Do the Refill Math
local delta_ms = math.max(0, now - last_ts)
local tokens_to_add = math.floor((delta_ms * capacity) / window_ms)
local current_tokens = math.min(capacity, last_tokens + tokens_to_add)

-- 3. The Decision
local allowed = 0
if current_tokens >= 1 then
    allowed = 1
    current_tokens = current_tokens - 1
    last_ts = now 
end

-- 4. Save and set auto-delete timer (Simple SET)
redis.call("SET", tokens_key, current_tokens, "PX", window_ms)
redis.call("SET", ts_key, last_ts, "PX", window_ms)

return allowed