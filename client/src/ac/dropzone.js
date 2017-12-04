import handleDataTransfer from './data-transfer';
import Alert from 'react-s-alert';

const unexpectedFilesMessage = 'Expected server_cfg.ini and entry_list.ini to be selected (or to be in the selected directory).';

const handleDrop = (cb) => {
  return (acceptedFiles, rejectedFiles, e) => {
    if (rejectedFiles.length > 0) {
      Alert.error(unexpectedFilesMessage);
      return;
    }

    switch (e.type) {
      case 'change': {
        const fd = new FormData();
        acceptedFiles.forEach((file) => {
          fd.append('configs', file);
        });
        cb(fd);
        break;
      }
      case 'drop': {
        handleDataTransfer(e, (e, formData, files) => {
          const validateFile = (entry) => {
            const filename = entry.replace(/^.*[\\/]/, '');
            return filename === 'server_cfg.ini' || filename === 'entry_list.ini';
          };

          if (!files || files.length !== 2 || files.filter(validateFile).length !== 2) {
            Alert.error(unexpectedFilesMessage);
            return;
          }

          cb(formData);
        });
        break;
      }
      default:
        Alert.error(`Unsupported event: ${e.type}`);
    }
  };
};

export default handleDrop;
