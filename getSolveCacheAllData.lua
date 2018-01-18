local qusetion_name_table = {}
local solve_team_table = {}
local cursor = 0
while true do
    local t_scan = redis.call('SCAN', cursor)

    for i=1,#t_scan[2] do
        qusetion_name_table[#qusetion_name_table + 1] = t_scan[2][i]

        local len = redis.call('LLEN', t_scan[2][i])
        local t_lrange = redis.call('LRANGE', t_scan[2][i], 0, len - 1)
        solve_team_table[#solve_team_table + 1] = t_lrange
    end
    -- to next
    cursor = t_scan[1]
    if t_scan[1] == '0' then
        break
    end
end

return {qusetion_name_table, solve_team_table}