-- local function array_add(m_arr, add_arr)
--     for i = 1, #add_arr then
--         m_arr[ #m_arr + 1] = add_arr[i]
--     end
--     return m_arr
-- end

local result = {}


local cursor = 0
local count = 0
while true do
    local t_scan = redis.call('SCAN', cursor)

    for i=1,#t_scan[2] do
        local r = {}
        r[#r + 1] = t_scan[2][i]
        r[#r + 1] = redis.call('HGET', t_scan[2][i], 'score')
        r[#r + 1] = redis.call('HGET', t_scan[2][i], 'update_time')

        result[#result + 1] = r
    end
    -- to next
    cursor = t_scan[1]
    count = count + #t_scan[2]
    if t_scan[1] == '0' then
        break
    end
end

return result