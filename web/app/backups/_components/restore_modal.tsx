import { useEffect, useState } from "react";
import {
  Modal, 
  ModalContent, 
  ModalHeader, 
  ModalBody, 
  ModalFooter, 
  Button, 
  useDisclosure, 
} from "@nextui-org/react";
import { NewBackupApi } from '@/components/api';
import { Backup, Restore, Status, Volume, TaskInput } from '@/components/spec';
import { TaskInputModal } from './task_input_modal'

export function RestoreBackup(props: {
  backup: Backup
}) {
  const {isOpen: optionsIsOpen, onOpen: optionsOnOpen, onOpenChange: optionsOnOpenChange} = useDisclosure();
  const {isOpen: progressIsOpen, onOpen: progressOnOpen, onOpenChange: progressOnOpenChange} = useDisclosure();

  return (<Button onPress={optionsOnOpen} isDisabled={props.backup.status === Status.Error}>
    Restore
    {optionsIsOpen && <TaskInputModal
      progressOnOpen={progressOnOpen}
      onOpenChange={optionsOnOpenChange}
      callApi={(taskInput: TaskInput, onClose: () => void) => {
        console.log(JSON.stringify(taskInput))
        var api = NewBackupApi()
        api.startRestore({
          id: props.backup.id,
          taskInput: {} as TaskInput,
        }).then(data => {
          onClose()
          progressOnOpen()
        }).catch(err => {
          alert("Error: " + err.error)
        })
      }}
    />}
    {progressIsOpen && <RestoreWaitModal
      id={props.backup.id}
      onOpenChange={progressOnOpenChange}
    />}
  </Button>);
}

function RestoreWaitModal(props: {
  id: string
  onOpenChange: () => void
}) {

  const [restore, setRestore] = useState<Restore>()
  const [error, setError] = useState<String>()

  useEffect(() => {
    const api = NewBackupApi()

    const id = setInterval(() => {
      api.getRestore({
        id: props.id
      }).then(data => {
        setRestore(restore)
        setError(undefined)

        if (data.status == Status.Error || data.status == Status.Completed) {
          clearInterval(id)
        }
      }).catch(err => {
        setError(err.response.text)
      })
    }, 1000)
    return () => clearInterval(id)
  }, [props.id])

  return (<Modal 
      isOpen={true} 
      placement="auto"
      onOpenChange={props.onOpenChange} 
    >
      <ModalContent>
        {(onClose) => (
          <>
            <ModalHeader className="flex flex-col gap-1">Waiting for Backup</ModalHeader>
            <ModalBody>
              {error && <pre>Error: {error}</pre>}
              {restore && <>
              <p>Id: {restore.id}</p>
              <p>Started: {JSON.stringify(restore.started)}</p>
              <p>Status: {restore.status}</p>
              {restore.error && <p>Error: {restore.error}</p>}
              {restore.volumes && <p>Volumes:
                <ul>
                  {restore.volumes.map((volume: Volume) => <li key={volume.name}>{JSON.stringify(volume)}</li>)}
                </ul>
              </p>}
              </>}
            </ModalBody>
            <ModalFooter>
              <Button color="danger" variant="light" onPress={onClose}>
                Cancel
              </Button>
            </ModalFooter>
          </>
        )}
      </ModalContent>
    </Modal>)
}