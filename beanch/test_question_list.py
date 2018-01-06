from lib.question import question
import config

q = question(config=config)

q_id_list = q.list()

print(q_id_list)