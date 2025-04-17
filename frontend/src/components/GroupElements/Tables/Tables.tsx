import * as React from 'react';
import './Tables.css'
import ContestTableRouting from './ContestTable/ContestTableRouting';
import UncheckedSubmissions from './UncheckedSubmissions/UncheckedSubmissions';


const Tables: React.FC = () => {
  return (
    <div className='tables'>
      <div className='contest-table'>
        <ContestTableRouting />
      </div>
      <UncheckedSubmissions />
    </div>
  );
};


export default Tables;
