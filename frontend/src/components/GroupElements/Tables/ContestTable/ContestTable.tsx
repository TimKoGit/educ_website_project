import * as React from 'react';
import { useNavigate } from 'react-router-dom';
import './ContestTable.css'


interface ContestInfo {
  id: string,
  name: string,
  time: string,
  duration: string
}


interface ContestInfoArray {
  contest_infos: ContestInfo[]
}


const ContestTable: React.FC<ContestInfoArray> = (contests: ContestInfoArray) => {
  const navigate = useNavigate();

  const handleRowClick = (contestId: string) => {
    navigate(`/contest/${contestId}`);
  };

  return (
    <div className='table-container contest-table-container'>
      <table>
        <thead>
          <tr>
            <th>Название</th>
            <th>Время старта</th>
            <th>Длительность</th>
          </tr>
        </thead>
        <tbody>
          {contests.contest_infos && contests.contest_infos.map((item) => (
            <tr
              key={item.id}
              onClick={() => handleRowClick(item.id)}
              className='clickable-row'
            >
              <td>{item.name}</td>
              <td>{item.time}</td>
              <td>{item.duration}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};


export default ContestTable;
