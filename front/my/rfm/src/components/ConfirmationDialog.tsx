import * as React from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';

interface ConfirmationDialogProps {
    open: boolean,
    title: string,
    content: string,
    handleClose: () => void,
    handleConfirm: () => void
}

const ConfirmationDialog: React.FC<ConfirmationDialogProps> = ({
    open,
    title,
    content,
    handleClose,
    handleConfirm
}) => {
    return (
        <Dialog
            open={open}
            onClose={handleClose}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description"
        >
            <DialogTitle id="alert-dialog-title">
                {title}
            </DialogTitle>
            <DialogContent>
                <DialogContentText id="alert-dialog-description">
                    {content}
                </DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose}>Cancelar</Button>
                <Button onClick={handleConfirm} autoFocus>
                    Confirmar
                </Button>
            </DialogActions>
        </Dialog>
    );
}

export default ConfirmationDialog;
