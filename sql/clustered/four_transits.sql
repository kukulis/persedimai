select  * -- TODO name columns instead
from clustered_arrival_travels32 c1
         join clustered_arrival_travels32 c2
              on c1.to_point = c2.from_point
                  and c1.arrival_cl = c2.departure_cl
         join clustered_arrival_travels32 c3
              on c2.to_point = c3.from_point
                  and c2.arrival_cl = c3.departure_cl
         join clustered_arrival_travels32 c4
              on c3.to_point = c4.from_point
                  and c3.arrival_cl = c4.departure_cl
where c1.from_point = '93456aaf-cf7e-4471-bfdb-3839145d7e73'
          and c4.to_point = 'fe2d35fe-698a-4f49-b38b-b08cbf61bfa0'
          and c4.arrival_cl >= 502536 && c4.arrival_cl <= 503280;