using DevSmtp.Core.Stores;

namespace DevSmtp.Core.Commands
{
    public sealed class SomlHandler : ICommandHandler<Soml, SomlResult>
    {
        private readonly IDataStore _dataStore;

        public SomlHandler(IDataStore dataStore)
        {
            this._dataStore = dataStore ?? throw new ArgumentNullException(nameof(dataStore));
        }

        public Task<SomlResult> ExecuteAsync(Soml command, CancellationToken cancellationToken = default)
        {
            throw new NotImplementedException();
        }
    }
}
