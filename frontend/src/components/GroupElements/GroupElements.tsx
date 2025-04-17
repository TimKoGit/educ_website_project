import * as React from 'react';
import { useEffect } from 'react';
import { Routes, Route, useNavigate, useParams } from 'react-router-dom';
import GroupSelection from './GroupSelection/GroupSelection'
import './GroupElements.css'
import Tables from './Tables/Tables';
import SubmissionForm from './SubmissionForm/SubmissionForm';
import CreateContestForm from './ContestForm/ContestForm';
import { useAppSelector, useAppDispatch } from '../../redux/hooks';
import AddTaskForm from './AddTaskForm/AddTaskForm';
import { fetchContestById } from '../../redux/slices/contestsSlice';
import { deleteGroup } from '../../redux/slices/authSlice';
import { deleteContest } from '../../redux/slices/contestsSlice';


const GroupElements: React.FC = () => {
  const user = useAppSelector((state) => state.auth.user);
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const { groupid, contestid } = useParams<{ groupid: string; contestid: string }>();

  const handleCreateContestClick = () => {
    navigate('/create-contest?groupid=' + groupid);
  };

  const contest = useAppSelector((state) => 
    contestid ? state.contests.contestsById[contestid] : null
  );
  const contestName = contest ? contest.name : 'Неизвестный контест';

  const { fetchContestLoading, fetchContestError } = useAppSelector((state) => (state.contests));

  useEffect(() => {
    if (contestid && !contest) {
      dispatch(fetchContestById(contestid));
    }
  }, [contestid, contest, dispatch]);

  const handleAddTaskClick = () => {
    navigate(`/add-task?contestid=${contestid}`);
  };

  const handleDeleteGroupClick = () => {
    if (groupid) {
      if (window.confirm('Вы уверены, что хотите удалить эту группу?')) {
        dispatch(deleteGroup(groupid))
          .then(() => {
            navigate('/');
          })
          .catch((error: any) => {
            console.error('Failed to delete group:', error);
            alert('Не удалось удалить группу. Попробуйте снова.');
          });
      }
    }
  };

  const handleDeleteContestClick = () => {
    if (contestid) {
      if (window.confirm('Вы уверены, что хотите удалить этот контест?')) {
        dispatch(deleteContest(contestid))
          .then(() => {
            navigate('/');
          })
          .catch((error: any) => {
            console.error('Failed to delete contest:', error);
            alert('Не удалось удалить контест. Попробуйте снова.');
          });
      }
    }
  };

  return (
    <div className='group-elements'>
      <GroupSelection />

      {contestid && (
        <div className='contest-name-section'>
          {fetchContestLoading ? (
            <p>Загрузка информации о контесте...</p>
          ) : fetchContestError ? (
            <p style={{ color: 'red' }}>{fetchContestError}</p>
          ) : (
            <h2>{contestName}</h2>
          )}
        </div>
      )}

      <Tables />
      <Routes>
        {user && user.role === 'teacher' && (
          <Route path='/create-contest' element={<CreateContestForm />} />
        )}
        {user && user.role === 'teacher' && (
          <Route path='submission/:submissionid' element={<SubmissionForm />} />
        )}
        {user && user.role === 'teacher' && (
          <Route path='/add-task' element={<AddTaskForm />} />
        )}
      </Routes>
      {user && user.role === 'teacher' && groupid && (
        <div className='button-group'>
          <div className='create-contest-section'>
            <button
              className='create-contest-btn'
              onClick={handleCreateContestClick}
            >
              Создать контест
            </button>
          </div>

          <div className='delete-group-section'>
            <button
              className='delete-group-btn'
              onClick={handleDeleteGroupClick}
            >
              Удалить группу
            </button>
          </div>
        </div>
      )}
      {user && user.role === 'teacher' && contestid && (
        <div className='button-group'>
          <div className='add-task-section'>
            <button
              className='add-task-btn'
              onClick={handleAddTaskClick}
            >
              Добавить задачу
            </button>
          </div>

          <div className='delete-contest-section'>
            <button
              className='delete-contest-btn'
              onClick={handleDeleteContestClick}
            >
              Удалить контест
            </button>
          </div>
        </div>
      )}
    </div>
  );
};


export default GroupElements;
