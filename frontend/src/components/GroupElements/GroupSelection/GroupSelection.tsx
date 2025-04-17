import * as React from 'react';
import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import './GroupSelection.css';
import { useAppSelector } from '../../../redux/hooks';
import GroupForm from './GroupForm/GroupForm';
import GroupJoinForm from './GroupJoinForm/GroupJoinForm';


const GroupSelection: React.FC = () => {

  const user = useAppSelector((state) => state.auth.user);
  const dbGroups = useAppSelector((state) => state.auth.groups);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [showJoinForm, setShowJoinForm] = useState(false);

  useEffect(() => {
    if (!user) {
      setShowCreateForm(false);
      setShowJoinForm(false);
    }
  }, [user]);

  const handleCreateGroupClick = () => {
    setShowCreateForm(!showCreateForm);
  };

  const handleJoinGroupClick = () => {
    setShowJoinForm(!showJoinForm);
  };

  return (
    <div className='blue-stripe'>
      <div className='left-side-outro'>
        <div className='left-side'>
          {dbGroups ? (dbGroups.map((group) => (
            <Link key={group.id} to={`/group/${group.id}`} className='left-button'>
              {group.name}
            </Link>
          ))) : (<div />)}
        </div>
      </div>

      <div className='right-side'>
        { user && (user.role === 'teacher') ? 
        (<button className='right-button' onClick={handleCreateGroupClick}>
          Создать новую группу
        </button>) : (<div />)}

        {showCreateForm && <GroupForm closeForm={() => setShowCreateForm(false)} />}

        { user && user.role === 'student' && (
          <>
            <button className='right-button' onClick={handleJoinGroupClick}>
              Присоединиться к группе
            </button>
            {showJoinForm && <GroupJoinForm closeForm={() => setShowJoinForm(false)} />}
          </>
        )}

      </div>
    </div>
  );
};


export default GroupSelection;
