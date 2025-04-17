import * as React from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import { RootState, AppDispatch } from '../../../redux/store';
import { acceptSubmission, declineSubmission, fetchSubmissionById, clearSpecificSubmission } from '../../../redux/slices/submissionsSlice';
import './SubmissionForm.css';


const SubmissionForm: React.FC = () => {
  const { submissionid } = useParams<{ submissionid: string }>();
  const dispatch = useDispatch<AppDispatch>();
  const navigate = useNavigate();
  const { specificSubmission, fetchSubmissionByIdLoading, fetchSubmissionByIdError } = useSelector((state: RootState) => state.submissions);
  const { user } = useSelector((state: RootState) => state.auth);

  React.useEffect(() => {
    if (submissionid && user?.id) {
      dispatch(fetchSubmissionById({submissionid}));
    }

    return () => {
      dispatch(clearSpecificSubmission());
    };
  }, [dispatch, submissionid, user]);

  const handleAccept = async () => {
    if (submissionid) {
      try {
        await dispatch(acceptSubmission(submissionid)).unwrap();
        navigate(-1);
      } catch (err: any) {
        console.error('Failed to accept submission:', err);
      }
    }
  };

  const handleDecline = async () => {
    if (submissionid) {
      try {
        await dispatch(declineSubmission(submissionid)).unwrap();
        navigate(-1);
      } catch (err: any) {
        console.error('Failed to decline submission:', err);
      }
    }
  };

  if (fetchSubmissionByIdLoading) {
    return <p>Загрузка посылки...</p>;
  }

  if (fetchSubmissionByIdError) {
    return <p className='error'>{fetchSubmissionByIdError}</p>;
  }

  if (!specificSubmission) {
    return <p>Посылка не найдена.</p>;
  }

  return (
    <div className='submission-form-container'>
      <h2>Информация о посылке</h2>
      <div className='submission-details'>
        <p><strong>Ученик:</strong> {specificSubmission.studentName}</p>
        <p><strong>Задача:</strong> {specificSubmission.taskName}</p>
        <div className='submission-photo'>
          <img src={specificSubmission.url} alt='Submission' />
        </div>
      </div>
      <div className='submission-actions'>
        <button className='accept-button' onClick={handleAccept}>Принять</button>
        <button className='decline-button' onClick={handleDecline}>Отклонить</button>
      </div>
    </div>
  );
};


export default SubmissionForm;
