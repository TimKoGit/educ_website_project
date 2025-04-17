import * as React from 'react';
import { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../../../redux/hooks';
import { addTask } from '../../../redux/slices/tasksSlice';
import { fetchContestById } from '../../../redux/slices/contestsSlice';
import './AddTaskForm.css';


const AddTaskForm: React.FC = () => {
  const [taskName, setTaskName] = useState('');
  const [taskUrl, setTaskUrl] = useState('');
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const location = useLocation();
  const searchParams = new URLSearchParams(location.search);
  const contestid = searchParams.get('contestid');
  const [error, setError] = useState<string | null>(null);
  const { addTaskLoading, addTaskError } = useAppSelector((state) => state.tasks);
  const contestsById = useAppSelector((state) => state.contests.contestsById);
  const contest = contestid ? contestsById[contestid] : null;
  const contestName = contest ? contest.name : 'Неизвестный контест';
  const { fetchContestLoading, fetchContestError } = useAppSelector((state) => state.contests);

  useEffect(() => {
    if (contestid && !contest) {
      dispatch(fetchContestById(contestid));
    }
  }, [contestid, contest, dispatch]);

  const handleInvalid = (e: React.FormEvent<HTMLInputElement>) => {
    e.preventDefault();
    if (!taskName) {
      setError('Введите название задачи');
    } else if (!taskUrl) {
      setError('Введите URL задачи');
    }
  };

  const onClose = () => {
    if (contestid) {
      navigate(`/contest/${contestid}`);
    } else {
      navigate('/');
    }
  };

  const handleCreateTask = async (e: React.FormEvent) => {
    e.preventDefault();
    if (taskName && taskUrl && contestid) {
      try {
        await dispatch(
          addTask({
            name: taskName,
            url: taskUrl,
            contestid
          })
        ).unwrap();
        onClose();
      } catch (err: any) {
        setError(err.message || 'Ошибка при добавлении задачи');
      }
    } else {
      alert('Все поля обязательны для заполнения.');
    }
  };

  return (
    <form onSubmit={handleCreateTask} className='add-task-form'>
      <h3 className='form-title'>Новая задача в контесте <strong>{fetchContestLoading ? 'Загрузка...' : contestName}</strong>!</h3>

      <div className='errors'>
        {error && <p style={{ color: 'red' }}>{error}</p>}
        {fetchContestError && <p style={{ color: 'red' }}>{fetchContestError}</p>}
        {addTaskError && <p style={{ color: 'red' }}>{addTaskError}</p>}
      </div>

      <div className='takNameInput'>
        <label htmlFor='taskName'>Название задачи:</label>
        <input
          type='text'
          id='taskName'
          value={taskName}
          onChange={(e) => setTaskName(e.target.value)}
          required
          onInvalid={(e) => handleInvalid(e)}
        />
      </div>

      <div className='taskUrlInput'>
        <label htmlFor='taskUrl'>URL задачи:</label>
        <input
          type='text'
          id='taskUrl'
          value={taskUrl}
          onChange={(e) => setTaskUrl(e.target.value)}
          required
          onInvalid={(e) => handleInvalid(e)}
        />
      </div>

      <div className='button-container'>
        <button type='submit' disabled={addTaskLoading}>
          {addTaskLoading ? 'Добавление...' : 'Добавить'}
        </button>
        <button type='button' onClick={onClose}>
          Назад
        </button>
      </div>
    </form>
  );
};


export default AddTaskForm;
