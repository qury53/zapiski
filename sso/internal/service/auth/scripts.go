package auth_service

import "github.com/redis/go-redis/v9"

var verifyEmailScript = redis.NewScript(`
-- KEYS[1] = code key
-- KEYS[2] = attempts key
-- ARGV[1] = hash отправленного кода (hex)
-- ARGV[2] = maxAttempts
-- ARGV[3] = TTL счётчика, мс
--
-- return: 1 ok | 0 wrong | -1 not found | -2 locked

local codeKey     = KEYS[1]
local attemptsKey = KEYS[2]
local submitted   = ARGV[1]
local maxAttempts = tonumber(ARGV[2])
local ttlMs       = tonumber(ARGV[3])

-- 1. код вообще есть?
local stored = redis.call('GET', codeKey)
if not stored then
  return -1
end

-- 2. страховка: счётчик уже на максимуме, а код почему-то жив
local current = tonumber(redis.call('GET', attemptsKey) or '0')
if current >= maxAttempts then
  redis.call('DEL', codeKey)
  return -2
end

-- 3. верный код: гасим и код, и счётчик (одноразовость)
if stored == submitted then
  redis.call('DEL', codeKey)
  redis.call('DEL', attemptsKey)
  return 1
end

-- 4. неверный: инкремент, TTL только на первой попытке
local n = redis.call('INCR', attemptsKey)
if n == 1 then
  local ttl = ttlMs
  local codeTTL = redis.call('PTTL', codeKey)  -- -1 = без TTL
  if codeTTL > 0 and codeTTL < ttl then
    ttl = codeTTL                              -- счётчик не переживёт код
  end
  redis.call('PEXPIRE', attemptsKey, ttl)
end

if n >= maxAttempts then
  redis.call('DEL', codeKey)                   -- блокировка = код уничтожен
  return -2
end

return 0
`)
var issueCodeScript = redis.NewScript(`
-- KEYS[1] = code | KEYS[2] = attempts | KEYS[3] = cooldown
-- ARGV[1] = hash кода | ARGV[2] = ttl кода, мс | ARGV[3] = cooldown, мс
-- return: 1 issued | 0 too soon

if redis.call('SET', KEYS[3], '1', 'PX', ARGV[3], 'NX') == false then
  return 0
end

redis.call('SET', KEYS[1], ARGV[1], 'PX', ARGV[2])
redis.call('DEL', KEYS[2])          -- сброс попыток при новой выдаче
return 1
`)
