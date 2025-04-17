import * as React from 'react';
import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../../../redux/hooks';
import { addNewContest } from '../../../redux/slices/contestsSlice';
import './ContestForm.css';


const CreateContestForm: React.FC = () => {
  const [contestName, setContestName] = useState('');
  const [startTime, setStartTime] = useState('');
  const [duration, setDuration] = useState('');
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const location = useLocation();
  const searchParams = new URLSearchParams(location.search);
  const groupid = searchParams.get('groupid');
  const [error, setError] = useState<string | null>(null);
  const { addNewContestLoading, addNewContestError } = useAppSelector((state) => state.contests);

  const handleInvalid = (e: React.FormEvent<HTMLInputElement>) => {
    e.preventDefault();
    if (!contestName) {
      setError('Введите название контеста');
    } else if (!startTime) {
      setError('Введите время начала');
    } else if (!duration) {
      setError('Введите продолжительность');
    }
  };

  const onClose = () => {
    if (groupid) {
      navigate(`/group/${groupid}`);
    } else {
      navigate('/');
    }
  };

  const group = useAppSelector((state) =>
    state.auth.groups?.find((g) => g.id === groupid)
  );

  const handleCreateContest = async (e: React.FormEvent) => {
    e.preventDefault();
    if (contestName && startTime && duration && groupid) {
      try {
        await dispatch(addNewContest({
          name: contestName,
          groupid,
          time: startTime,
          duration,
        })).unwrap();
        onClose();
      } catch (err: any) {
        setError(err.message || 'Ошибка при создании контеста');
      }
    } else {
      alert('Все поля обязательны для заполнения.');
    }
  };

  return (
    <form onSubmit={handleCreateContest} className='create-contest-form'>
      <h3 className='form-title'>Новый контест в группе <strong>{group?.name || 'неизвестная группа'}</strong>!</h3>

      {error && !addNewContestError  && <p style={{ color: 'red' }}>{error}</p>}
      {addNewContestError && <p style={{ color: 'red' }}>{addNewContestError}</p>}

      <label htmlFor='contestName'>Название контеста:</label>
      <input
        type='text'
        id='contestName'
        value={contestName}
        onChange={(e) => setContestName(e.target.value)}
        required
        onInvalid={(e) => handleInvalid(e)}
      />

      <label htmlFor='startTime'>Время начала:</label>
      <input
        type='datetime-local'
        id='startTime'
        value={startTime}
        onChange={(e) => setStartTime(e.target.value)}
        required
        onInvalid={(e) => handleInvalid(e)}
      />

      <label htmlFor='duration'>Продолжительность (в часах):</label>
      <input
        type='number'
        id='duration'
        value={duration}
        onChange={(e) => setDuration(e.target.value)}
        required
        onInvalid={(e) => handleInvalid(e)}
      />

      <div className='button-container'>
        <button type='submit' disabled={addNewContestLoading}>
          {addNewContestLoading ? 'Создание...' : 'Создать'}
        </button>
        <button type='button' onClick={onClose}>
          Назад
        </button>
      </div>
    </form>
  );
};

export default CreateContestForm;
