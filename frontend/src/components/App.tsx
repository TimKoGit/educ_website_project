import * as React from 'react';
import { Route, BrowserRouter as Router, Routes} from 'react-router-dom'
import UserMenu from './UserMenu/UserMenu';
import GroupElements from './GroupElements/GroupElements';
import './App.css'


export default function App(): React.JSX.Element {
  return (
    <Router>
      <div className='app'>
        <UserMenu />
        <Routes>
          <Route path='/group/:groupid/*' element={<GroupElements />} />
          <Route path='/contest/:contestid/*' element={<GroupElements />} />
          <Route path='*' element={<GroupElements/>}/>
        </Routes>
      </div>
    </Router>
  )
}
