import { useState } from "react";
import {
  Selection,
  Listbox,
  Checkbox,
  ListboxItem,
  ListboxSection, 
  Modal, 
  ModalContent, 
  ModalHeader, 
  ModalBody, 
  ModalFooter, 
  Button, 
  Input
} from "@nextui-org/react";
import { NewBackupApi } from '@/components/api';
import { Backup, TaskInput, Filters } from '@/components/spec';

export function TaskInputModal(props: {
  progressOnOpen: () => void
  onOpenChange: () => void
  callApi: (taskInput: TaskInput, onClose: () => void) => void
}) {
  
  var filters = new Filters([], [], [], [])
  var taskInput = new TaskInput(filters)

  const [selectedKeys, setSelectedKeys] = useState(new Set(["labelled", "writers", "linked"]));

  return (<Modal 
      isOpen={true} 
      placement="auto"
      onOpenChange={props.onOpenChange} 
    >
      <ModalContent>
        {(onClose) => (
          <>
            <ModalHeader className="flex flex-col gap-1">Create Backup</ModalHeader>
            <ModalBody>
              <Checkbox radius="full" onChange={(e) => taskInput.no_dry_run = e.target.checked}>
                No Dry Run
              </Checkbox>
              
              <ListInput label="Include Names" updated={(s) => filters.include_names = s}/>
              <ListInput label="Include Volumes" updated={(s) => filters.include_volumes = s}/>
              <ListInput label="Exclude Names" updated={(s) => filters.exclude_names = s}/>
              <ListInput label="Exclude Volumes" updated={(s) => filters.exclude_volumes = s}/>

              <Listbox 
                aria-label="Stop Modes"
                disallowEmptySelection
                selectionMode="multiple"
                selectedKeys={selectedKeys}
                onSelectionChange={(keys: Selection) => {
                  setSelectedKeys(keys as Set<string>)
                  taskInput.stop_modes = keys
                }}
              >
                <ListboxSection title="Stop Modes"> 
                  <ListboxItem key="all">All</ListboxItem>
                  <ListboxItem key="labelled">Labelled</ListboxItem>
                  <ListboxItem key="writers">Writers</ListboxItem>
                  <ListboxItem key="attached">Attached</ListboxItem>
                  <ListboxItem key="linked">Linked</ListboxItem>
                </ListboxSection>
              </Listbox>

              <Input
                label="Label Prefix"
                defaultValue="bckupr"
                onChange={(e) => taskInput.label_prefix = e.target.value}
              />
            </ModalBody>
            <ModalFooter>
              <Button color="danger" variant="light" onPress={onClose}>
                Cancel
              </Button>
              <Button color="primary" onPress={() => props.callApi(taskInput, onClose)}>
                Start
              </Button>
              
            </ModalFooter>
          </>
        )}
      </ModalContent>
    </Modal>)
}

function ListInput(props: {
  label: string
  updated: (items: string[]) => void
}) {
  const [inputted, setInputted] = useState("")
  const [inputs, setInputs] = useState(new Set<string>())

  return (<>
    <div className="flex">
      <Input isClearable
        label={props.label}
        placeholder="..."
        value={inputted}
        onClear={() => setInputted("")}
        onChange={(e) => setInputted(e.target.value)}
      />
      <div className="content-around">
        <Button className="h-full min-w" onClick={() => {
          inputs.add(inputted)
          setInputted("")
          const updated = new Set(inputs)
          setInputs(updated)
          props.updated(Array.from(updated))
        }}>+</Button>
      </div>
    </div>
    {Array.from(inputs).map((name: string) => <div key={name} className="flex">
      <Input defaultValue={name} isReadOnly/>
      <div className="content-around">
        <Button className="h-full min-w" onClick={() => {
          inputs.delete(name)
          const updated = new Set(inputs)
          setInputs(updated)
          props.updated(Array.from(updated))
        }}>-</Button>
      </div>
    </div>)}
  </>)
}
