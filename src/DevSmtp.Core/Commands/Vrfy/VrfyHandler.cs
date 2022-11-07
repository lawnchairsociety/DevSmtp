using DevSmtp.Core.Stores;

namespace DevSmtp.Core.Commands
{
    public sealed class VrfyHandler : ICommandHandler<Vrfy, VrfyResult>
    {
        private readonly IDataStore _dataStore;

        public VrfyHandler(IDataStore dataStore)
        {
            this._dataStore = dataStore ?? throw new ArgumentNullException(nameof(dataStore));
        }

        public Task<VrfyResult> ExecuteAsync(Vrfy command, CancellationToken cancellationToken = default)
        {
            throw new NotImplementedException();
        }
    }
}
