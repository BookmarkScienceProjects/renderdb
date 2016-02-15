var scene
var camera
var headlight
var controls
var renderer

init()
loadData()
animate()

function init () {
  scene = new THREE.Scene()

  camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 1, 1000)
  scene.add(camera)

  controls = new THREE.TrackballControls(camera)
  controls.rotateSpeed = 1.0
  controls.zoomSpeed = 1.2
  controls.panSpeed = 0.8
  controls.noZoom = false
  controls.noPan = false
  controls.staticMoving = true
  controls.dynamicDampingFactor = 0.3
  controls.keys = [ 65, 83, 68 ]
  controls.addEventListener('change', render)

  controls.position = new THREE.Vector3(10, -9, 28)
  controls.target = new THREE.Vector3(20, -5, 35)

  var ambient = new THREE.AmbientLight(0x202020)
  scene.add(ambient)

  headlight = new THREE.DirectionalLight(0xffffff, 0.7)
  scene.add(headlight)

  renderer = new THREE.WebGLRenderer()
  renderer.setPixelRatio(window.devicePixelRatio)
  renderer.setSize(window.innerWidth, window.innerHeight)
  renderer.setClearColor(0x100000)

  document.body.appendChild(renderer.domElement)
}

function loadData () {
  $.post({
    url: '/geometry/view',
    data: JSON.stringify(
      {
        bounds: {
          min: [12, -9, 29],
          max: [28, 1.8, 40]
        },
        eyePosition: [10, -9, 28]
    })}).done(function (data) {
    data.forEach(function (id) {
      $.get('/geometry/' + id).done(function (entry) {
        var loader = new THREE.OBJLoader()
        data = atob(entry.geometryData)
        obj = loader.parse(data)
        scene.add(obj)
        
        requestAnimationFrame(render)
      })
    })
    camera.position = new THREE.Vector3(10, -9, 28)
    camera.lookAt(new THREE.Vector3(20, -5, 35))
  })
}

function animate () {
  requestAnimationFrame(animate)
  controls.update()
}

function render () {
  headlight.position.copy(camera.position)
  renderer.render(scene, camera)
}
