SELECT A.SHIP_ID,
       REPLACE(REGEXP_SUBSTR(OP_PICURL, '[^,]+', 1, 1), 'group', 'https://file.yuanfusc.com/group') AS url1,
       REPLACE(REGEXP_SUBSTR(OP_PICURL, '[^,]+', 1, 2), 'group', 'https://file.yuanfusc.com/group') AS url2,
       REPLACE(REGEXP_SUBSTR(OP_PICURL, '[^,]+', 1, 3), 'group', 'https://file.yuanfusc.com/group') AS url3,
       REPLACE(REGEXP_SUBSTR(OP_PICURL, '[^,]+', 1, 4), 'group', 'https://file.yuanfusc.com/group') AS url4,
       REPLACE(REGEXP_SUBSTR(OP_PICURL, '[^,]+', 1, 5), 'group', 'https://file.yuanfusc.com/group') AS url5
FROM TMS_APP_TASK_ORDER_DETAIL B
         JOIN TMS_APP_TASK_ORDER A ON A.TASK_ID = B.TASK_ID
WHERE B.OP_CODE = '004'
  AND A.SHIP_ID IN (:keyword);
